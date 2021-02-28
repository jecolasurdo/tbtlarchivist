package analyst

import (
	"context"
	"io"
)

// A CommandBuilder is able to build commands.
type CommandBuilder interface {
	CommandContext(context.Context, string, ...string) Command
}

// A Command is able to Start a command, Wait for it to complete, and
// communicate with it bia stdin and stdout.
type Command interface {
	StdoutPipe() (ReadCloser, error)
	StdinPipe() (WriteCloser, error)
	StderrPipe() (ReadCloser, error)
	Start() error
	Wait() error
}

// ReadCloser is an io.ReadCloser that has been reimplemented to ease mock
// generation.
type ReadCloser interface {
	io.ReadCloser
}

// WriteCloser is an io.ReadCloser that has been reimplemented to ease mock
// generation.
type WriteCloser interface {
	io.WriteCloser
}
