package utils

import "context"

// ContextIsDone returns true if the supplied context is reporting that it is
// done.
func ContextIsDone(ctx context.Context) bool {
	return ctx.Err() != nil
}
