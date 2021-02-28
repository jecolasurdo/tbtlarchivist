package utils

import (
	"context"
	"fmt"
	"time"
)

type BackoffAPI interface {
	Wait() error
	Reset()
}

// LinearBackoff provides a means of linearly increasing a wait period after each
// subsequent call to a Wait method.
type LinearBackoff struct {
	ctx               context.Context
	maxCumulativeWait time.Duration
	currentWait       time.Duration
	increment         time.Duration
	cumulativeWait    time.Duration
}

// NewLinearBackoff initializes a backoff. increment defines the amount to increase
// the wait duration after each call to Wait. maxCumulativeWait is the maximum
// amount of total time that Backoff will wait across consecutive calls to
// Wait before Wait will return an error.
func NewLinearBackoff(ctx context.Context, increment, maxCumulativeWait time.Duration) *LinearBackoff {
	return &LinearBackoff{
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
func (lb *LinearBackoff) Wait() error {
	if lb.cumulativeWait > lb.maxCumulativeWait {
		return fmt.Errorf("LinearBackoff: wait would exceeed maximum cumulative wait duration")
	}

	// Initial wait is zero, we accumulate after waiting, not before.
	select {
	case <-lb.ctx.Done():
		return fmt.Errorf("LinearBackoff: parent context done")
	case <-time.After(lb.currentWait):
	}

	lb.currentWait += lb.increment
	lb.cumulativeWait += lb.currentWait

	return nil
}

// Reset returns LinearBackoff to its initial state. Restoring the cumulative wait
// time to zero, and the current wait period back to zero.
func (lb *LinearBackoff) Reset() {
	lb.currentWait = 0
	lb.cumulativeWait = 0
}

type ConstantBackoff struct {
	ctx               context.Context
	rate              time.Duration
	maxCumulativeWait time.Duration
	cumulativeWait    time.Duration
}

func NewConstantBackoff(ctx context.Context, rate, maxCumulativeWait time.Duration) *ConstantBackoff {
	return &ConstantBackoff{
		ctx:               ctx,
		rate:              rate,
		maxCumulativeWait: maxCumulativeWait,
	}
}

func (cb *ConstantBackoff) Wait() error {
	if cb.cumulativeWait > cb.maxCumulativeWait {
		return fmt.Errorf("ConstantBackoff: wait would exceeed maximum cumulative wait duration")
	}

	select {
	case <-cb.ctx.Done():
		return fmt.Errorf("ConstantBackoff: parent context done")
	case <-time.After(cb.rate):
	}

	cb.cumulativeWait += cb.rate

	return nil
}

func (cb *ConstantBackoff) Reset() {
	cb.cumulativeWait = 0
}
