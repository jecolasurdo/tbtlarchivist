package rustanalyst

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/analyst"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/utils"
	"google.golang.org/protobuf/proto"
)

// The Adapter spawns a child analyst-rust process, and marshals messages
// between the caller and the child process.
type Adapter struct {
	errorSource         chan (error)
	completedItemSource chan (*contracts.CompletedResearchItem)
	done                chan (struct{})

	CmdBuilder analyst.CommandBuilder

	// PathResolver is a function that returns the path to the analyst process
	// to be spawned. If this value is nil, DefaultPathResolver is used.
	PathResolver func() (string, error)
}

// DefaultPathResolver is used to locate an "analyzerd" binary in the local
// PATH.
func DefaultPathResolver() (string, error) {
	path, err := exec.LookPath("analyzerd")
	if err != nil {
		return "", fmt.Errorf("error occured while locating the analyzer\n%v", err)
	}
	return path, nil
}

// Run starts a rust analyst as a child process, pipes pendingResearch to the
// process via stdin, and listens for results on stdout. The returned
// CompletedResearchItem and error channels will remain open until all work is
// completed, at which time they are both closed.  Any errors that occur before
// the adapter begins polling stdout will result in the closure of the
// CompletedResearchItem and error channels. However, any errors that occur
// while processing stdout are streamed to the outbound error channel. The
// adapter will continue to poll stdout until the pipe is either closed or the
// parent context is cancelled. If the parent context reports that it is Done,
// the child process is immediately killed (SIGKILL is sent to the child
// process).
func (a *Adapter) Run(ctx context.Context, pendingResearch *contracts.PendingResearchItem) {
	if a.CmdBuilder == nil {
		a.CmdBuilder = new(analyst.ExecFacade)
	}

	if a.PathResolver == nil {
		a.PathResolver = DefaultPathResolver
	}

	a.completedItemSource = make(chan *contracts.CompletedResearchItem)
	a.errorSource = make(chan error)
	a.done = make(chan struct{})
	go func() {
		defer close(a.completedItemSource)
		defer close(a.errorSource)
		defer close(a.done)

		path, err := a.PathResolver()
		if err != nil {
			a.errorSource <- err
			return
		}

		innerCtx, cancel := context.WithCancel(ctx)
		cmd := a.CmdBuilder.CommandContext(innerCtx, path)
		defer cancel()

		stderr, err := cmd.StderrPipe()
		if err != nil {
			a.errorSource <- err
			return
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			a.errorSource <- err
			return
		}

		stdin, err := cmd.StdinPipe()
		if err != nil {
			a.errorSource <- err
			return
		}

		err = cmd.Start()
		if err != nil {
			a.errorSource <- err
			return
		}

		pendingBytes, err := proto.Marshal(pendingResearch)
		if err != nil {
			a.errorSource <- err
			return
		}

		_, writeErr := stdin.Write(pendingBytes)
		if writeErr != nil {
			a.errorSource <- writeErr
		}

		closeErr := stdin.Close()
		if closeErr != nil {
			a.errorSource <- closeErr
		}

		if writeErr != nil || closeErr != nil {
			return
		}

		stdoutBackoff := utils.NewLinearBackoff(ctx, 100*time.Millisecond, 10*time.Second)
		stdoutScanner := utils.NewFrameScanner(stdout, stdoutBackoff)
		recordSource := stdoutScanner.Poll()

		stderrBackoff := utils.NewConstantBackoff(ctx, 1*time.Second, 24*365*time.Hour)
		stderrScanner := utils.NewFrameScanner(stderr, stderrBackoff)
		stderrSource := stderrScanner.Poll()

	loop:
		for {
			select {
			case <-ctx.Done():
				a.errorSource <- ctx.Err()
				break loop
			case record, open := <-recordSource:
				if !open {
					break loop
				}
				completedResearchItem := new(contracts.CompletedResearchItem)
				err = proto.Unmarshal(record, completedResearchItem)
				if err != nil {
					a.errorSource <- err
				} else {
					a.completedItemSource <- completedResearchItem
				}
			case err, open := <-stderrSource:
				if !open {
					break loop
				}
				if err != nil {
					a.errorSource <- fmt.Errorf(string(err))
				}
			}
		}

		if stdoutScanner.Err() != nil {
			a.errorSource <- stdoutScanner.Err()
		}

		if stderrScanner.Err() != nil {
			a.errorSource <- stderrScanner.Err()
		}

		a.errorSource <- cmd.Wait()
	}()
}

// Errors provides access to errors that are produced after Run called.
func (a *Adapter) Errors() <-chan (error) {
	return a.errorSource
}

// CompletedWorkItems provides access to a stream of completed work items.
func (a *Adapter) CompletedWorkItems() <-chan *contracts.CompletedResearchItem {
	return a.completedItemSource
}

// Done returns a channel that blocks until the adapter is done running.
func (a *Adapter) Done() <-chan (struct{}) {
	return a.done
}

var _ analyst.Analyzer = (*Adapter)(nil)
