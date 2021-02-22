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
	StdoutPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	Start() error
	Wait() error
}
