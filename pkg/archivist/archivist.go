package archivist

import (
	"context"
	"math/rand"
	"runtime"
	"time"
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

// An ArchiveWorker is anything that 1) provides a means of initializing a
// stream of data from some source, and 2) provides a method that can be used
// by a consumer to process each datum from that stream.
type ArchiveWorker interface {
	InitializeDataStream(context.Context) (<-chan interface{}, error)
	ProcessDatum(context.Context, interface{}) error
}

// Initialize activates a set of archive-workers. As each worker is
// initialized, it returns a stream of data to be processed. These streams are
// then forwarded to a poller.  The poller selects one of the streams at
// random.  If the stream has data, one datum is dequeued and sent to the
// appropriate worker's ProcessDatum method. If the stream has no data ready,
// the poller selects another worker's stream at random, and process repeats.
// The poll continues so long as all streams remain open. The general status of
// the API can be monitored via the API.Errors channel.  API.Errors returns a
// stream of any errors that might arise during operation.  The API.Errors
// channel remains open until all workers have safely wound down.  Thus,
// API.Errors can/should be used by the caller as a waiter to avoid premature
// termination of an application.  parentCtx is propogated to all downstream
// workers, and should be used to safely broadcast cancellation requests to the
// API.
func Initialize(parentCtx context.Context, workers []ArchiveWorker) *API {
	ctx, cancel := context.WithCancel(parentCtx)
	a := new(API)

	errSrc := make(chan error)
	a.Errors = errSrc

	go func() {
		defer close(errSrc)

		dataStreams := make([]<-chan interface{}, len(workers))
		dataProcessors := make([]func(context.Context, interface{}) error, len(workers))
		for i, worker := range workers {
			dataStream, err := worker.InitializeDataStream(ctx)
			dataStreams[i] = dataStream
			dataProcessors[i] = worker.ProcessDatum
			if err != nil {
				errSrc <- err
				cancel()
			}
		}

		rand.Seed(time.Now().UnixNano())
		openChannelCount := len(workers)
		for openChannelCount > 0 {
			i := rand.Intn(len(workers))
			select {
			case item, ok := <-dataStreams[i]:
				if !ok {
					openChannelCount--
				} else {
					if err := dataProcessors[i](ctx, item); err != nil {
						errSrc <- err
					}
				}
			default:
				// If none of the channels we're polling have any work ready,
				// we can end up in a busy loop. We call Gosched to prevent
				// unintentionally hogging the processor when there's no work
				// to do.
				runtime.Gosched()
			}
		}
	}()

	return a
}
