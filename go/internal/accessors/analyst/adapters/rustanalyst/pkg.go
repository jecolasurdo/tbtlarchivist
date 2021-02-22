package rustanalyst

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/analyst"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/utils"
	"google.golang.org/protobuf/proto"
)

// The Adapter spawns a child analyst-rust process, and marshals messages
// between the caller and the child process.
type Adapter struct {

	// PathResolver is a function that returns the path to the analyst process
	// to be spawned. If this value is nil, DefaultPathResolver is used.
	PathResolver func() (string, error)
}

// DefaultPathResolver is used to locate an "analyzer" binary in the local
// PATH.
func DefaultPathResolver() (string, error) {
	path, err := exec.LookPath("analyzer")
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
func (a *Adapter) Run(ctx context.Context, pendingResearch *contracts.PendingResearchItem) (<-chan *contracts.CompletedResearchItem, <-chan error) {
	if a.PathResolver == nil {
		a.PathResolver = DefaultPathResolver
	}

	completedItemSource := make(chan *contracts.CompletedResearchItem)
	errorSource := make(chan error)
	go func() {
		defer close(completedItemSource)
		defer close(errorSource)

		path, err := a.PathResolver()
		if err != nil {
			errorSource <- err
			return
		}

		cmd := exec.Command(path)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errorSource <- err
			return
		}

		stdin, err := cmd.StdinPipe()
		if err != nil {
			errorSource <- err
			return
		}

		err = cmd.Start()
		if err != nil {
			errorSource <- err
			return
		}

		pendingBytes, err := proto.Marshal(pendingResearch)
		if err != nil {
			errorSource <- err
			return
		}

		_, writeErr := stdin.Write(pendingBytes)
		if writeErr != nil {
			errorSource <- writeErr
		}

		closeErr := stdin.Close()
		if closeErr != nil {
			errorSource <- closeErr
		}

		if writeErr != nil || closeErr != nil {
			errorSource <- cmd.Process.Kill()
			return
		}

		scanner := bufio.NewScanner(stdout)
		frameScanner := new(utils.FrameScanner)
		scanner.Split(frameScanner.ScanFrames)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				errorSource <- ctx.Err()
				errorSource <- cmd.Process.Kill()
				return
			default:
			}

			if scanner.Err() != nil {
				errorSource <- scanner.Err()
			}
			completedResearchItem := new(contracts.CompletedResearchItem)
			err = proto.Unmarshal(scanner.Bytes(), completedResearchItem)
			if err != nil {
				errorSource <- err
			}
			completedItemSource <- completedResearchItem
		}

		errorSource <- cmd.Wait()
	}()

	return completedItemSource, errorSource
}

var _ analyst.Analyzer = (*Adapter)(nil)
