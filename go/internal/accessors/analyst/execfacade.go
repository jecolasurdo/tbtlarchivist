package analyst

import (
	"context"
	"os/exec"
)

// ExecFacade is a facade for the exec package.
type ExecFacade struct{}

// CommandContext is a facade for exec.CommandContext
func (*ExecFacade) CommandContext(ctx context.Context, name string, arg ...string) Command {
	return &ExecCmdFacade{
		cmd: exec.CommandContext(ctx, name, arg...),
	}
}

// ExecCmdFacade is a facade for an exec.Cmd
type ExecCmdFacade struct {
	cmd *exec.Cmd
}

// StdoutPipe is a facade for exec.Cmd.StdoutPipe.
func (c *ExecCmdFacade) StdoutPipe() (ReadCloser, error) {
	return c.cmd.StdoutPipe()
}

// StdinPipe is a facade for exec.Cmd.StdinPipe.
func (c *ExecCmdFacade) StdinPipe() (WriteCloser, error) {
	return c.cmd.StdinPipe()
}

// Start is a facade for exec.Cmd.Start.
func (c *ExecCmdFacade) Start() error {
	return c.cmd.Start()
}

// Wait is a facade for exec.Cmd.Wait.
func (c *ExecCmdFacade) Wait() error {
	return c.cmd.Wait()
}
