package utils

import (
	"context"
	"fmt"
	"time"
)

// Backoff provides a means of linearly increasing a wait period after each
// subsequent call to a Wait method.
type Backoff struct {
	ctx               context.Context
	maxCumulativeWait time.Duration
	currentWait       time.Duration
	increment         time.Duration
	cumulativeWait    time.Duration
}

// NewBackoff initializes a backoff. increment defines the amount to increase
// the wait duration after each call to Wait. maxCumulativeWait is the maximum
// amount of total time that Backoff will wait across consecutive calls to
// Wait before Wait will return an error.
func NewBackoff(ctx context.Context, increment, maxCumulativeWait time.Duration) *Backoff {
	return &Backoff{
		ctx:               ctx,
		maxCumulativeWait: maxCumulativeWait,
		increment:         increment,
	}
}

// Wait blocks for a period of time. The duration of the wait period increases
// linearly with each call to Wait. The initial blocking period is always zero,
// so only the second and subsequent calls to Wait effectively block.  If the
// total time spent waiting would exceed the maximum cumulative wait duration,
// the method does not block, and instead immediately returns an error.
// Context cancellation is given higher priority than the current wait period.
// Thus, if the parent context is cancelled before the current wait period is
// complete, the method will immediately return a "parent context done" error.
// When Wait returns an error, it should not be called again unless the Reset
// method is called first. Calling the Wait method after it has returned an
// error but without first calling Reset is undefined, and may produce
// undesirable results.
func (b *Backoff) Wait() error {
	if b.cumulativeWait > b.maxCumulativeWait {
		return fmt.Errorf("Backoff: wait would exceeed maximum cumulative wait duration")
	}

	// Initial wait is zero, we accumulate after waiting, not before.
	select {
	case <-b.ctx.Done():
		return fmt.Errorf("Backoff: parent context done")
	case <-time.After(b.currentWait):
	}

	b.currentWait += b.increment
	b.cumulativeWait += b.currentWait

	return nil
}

// Reset returns Backoff to its initial state. Restoring the cumulative wait
// time to zero, and the current wait period back to zero.
func (b *Backoff) Reset() {
	b.currentWait = 0
	b.cumulativeWait = 0
}
