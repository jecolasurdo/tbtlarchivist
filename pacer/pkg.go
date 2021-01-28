// Package pacer provides capabilities for throttling the rate at which a
// channel is flushed.
package pacer

import (
	"context"
	"runtime"
)

// Queue .
type Queue struct {
	backlog chan func() error
}

// Enqueue places as action in the queue where it will be wait to be called.
func (q Queue) Enqueue(action func() error) {
	if q.backlog == nil {
		q.backlog = make(chan func() error)
	}
	q.backlog <- action
}

// Poll initializes the internall poller, which pumps the queue, executing
// queued actions in a FIFO order. If the supplied context is cancelled:
// - the poller is immediately stopped
// - its error channel is closed
// - the underlaying queue is immediately closed and deallocated
// - any remaining work in the queue will be lost as soon as ctx is
//   cancelled.
// It is only safe to call Poll a second time if the context for the previous
// call has been cancelled and is done.
func (q Queue) Poll(ctx context.Context) <-chan error {
	if q.backlog == nil {
		q.backlog = make(chan func() error)
	}
	errSource := make(chan error)
	go func() {
		defer close(errSource)
		for {
			select {
			case <-ctx.Done():
				errSource <- ctx.Err()
				q.backlog = nil
				return
			case action := <-q.backlog:
				errSource <- action()
			default:
				runtime.Gosched()
			}
		}
	}()
	return errSource
}
