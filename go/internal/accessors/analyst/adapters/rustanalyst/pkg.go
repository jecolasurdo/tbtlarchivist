package rustanalyst

import (
	"bufio"
	"context"
	"os/exec"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/analyst"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
	"google.golang.org/protobuf/proto"
)

// The Adapter spawns a child analyst-rust process, and marshals messages between
// the caller and the child processes.
type Adapter struct{}

// Run starts a rust analyst as a child process, pipes pendingResearch to the
// process via stdin, and listens for results on stdout.
func (a *Adapter) Run(ctx context.Context, pendingResearch *contracts.PendingResearchItem) (<-chan *contracts.CompletedResearchItem, <-chan error) {
	completedItemSource := make(chan *contracts.CompletedResearchItem)
	errorSource := make(chan error)
	go func() {
		defer close(completedItemSource)
		defer close(errorSource)

		path, err := exec.LookPath("analyzer")
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
		go func() {
			for scanner.Scan() {

			}
		}()

	}()

	return completedItemSource, errorSource
}

var _ analyst.Analyzer = (*Adapter)(nil)
