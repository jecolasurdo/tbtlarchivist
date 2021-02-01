package archivist

import (
	"context"
)

// API is an instance of an archivist. An archivist is responsible for the
// following tasks:
//  - Recording episode metadata as reported from the curators.
//  - Recording clip metadata as reported from the curators.
//  - Creating work for researchers.
//  - Checking in work returned from the researchers.
type API struct {
	Errors <-chan error
}

type worker struct {
	source   <-chan interface{}
	delegate func(interface{}) error
}

// Initialize activates a set of workers. Each worker returns a stream of data
// to be processed. The streams are processed across the workers as a
// round-robin.  The general status of the API can be monitored via the
// API.Errors channel.  API.Errors returns a stream of any errors that might
// arise during operation.  The API.Errors channel remains open until all
// workers have safely wound down.  Thus, API.Errors can/should be used by the
// caller as a waiter to avoid premature termination of an application.
// parentCtx is propogated to all downstream workers, and should be used to
// safely broadcast cancellation requests to the API.
func Initialize(parentCtx context.Context) *API {
	const numberOfSourceChannels = 4
	ctx, cancel := context.WithCancel(parentCtx)
	a := new(API)

	errSrc := make(chan error)
	a.Errors = errSrc

	go func() {
		defer close(errSrc)

		// Each of the following blocks calls `cancel()` if an error is
		// encountered.  If a subsequent block encounters an error, it
		// broadcasts a cancellation upon the receipt of which, each preceding
		// block can wind down and close their channels cleanly. The polling
		// loop keeps track of how many channels are open. Once all channels
		// have closed, the loop is stopped, and the error channel is closed.
		// This signals to upstream consumers who are monitoring the API.Errors
		// channel, that everything has wound down cleanly.
		ces, err := a.getCuratedEpisodeSource(ctx)
		if err != nil {
			cancel()
			errSrc <- err
		}

		ccs, err := a.getCuratedClipSource(ctx)
		if err != nil {
			cancel()
			errSrc <- err
		}

		prs, err := a.getPendingResearchSource(ctx)
		if err != nil {
			cancel()
			errSrc <- err
		}

		crs, err := a.getCompletedResearchSource(ctx)
		if err != nil {
			cancel()
			errSrc <- err
		}

		openChannelCount := numberOfSourceChannels
		for openChannelCount > 0 {
			select {
			case ce, ok := <-ces:
				if !ok {
					openChannelCount--
				} else {
					catch(a.processCuratedEpisode(ctx, ce), errSrc)
				}
			case cc, ok := <-ccs:
				if !ok {
					openChannelCount--
				} else {
					catch(a.processCuratedClip(ctx, cc), errSrc)
				}
			case pr, ok := <-prs:
				if !ok {
					openChannelCount--
				} else {
					catch(a.processPendingResearch(ctx, pr), errSrc)
				}
			case cr, ok := <-crs:
				if !ok {
					openChannelCount--
				} else {
					catch(a.processCompletedResearch(ctx, cr), errSrc)
				}
			}
		}
	}()

	return a
}

func catch(err error, ch chan<- error) {
	if err != nil {
		ch <- err
	}
}
