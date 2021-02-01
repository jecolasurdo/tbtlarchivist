package archivist

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
	"github.com/jecolasurdo/tbtlarchivist/pkg/messagebus/messagebusiface"
)

// A CuratedEpisodeWorker can initialize a stream of inbound curated episodes,
// and supplies a method for processing each episode.
type CuratedEpisodeWorker struct{}

// InitializeDataStream opens a stream of inbound episode metadata that needs
// to be processed.
func (c *CuratedEpisodeWorker) InitializeDataStream(ctx context.Context, msgBus messagebusiface.MessageBus) (<-chan *messagebusiface.MessageBusMessage, error) {
	episodeSource := make(chan *messagebusiface.MessageBusMessage)
	go func() {
		defer close(episodeSource)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			msg := msgBus.Receive()
			if msg != nil {
				episodeSource <- msg
			} else {
				runtime.Gosched()
			}
		}
	}()
	return episodeSource, nil
}

// ProcessDatum processes an episode, determing whether or not it should be
// added to the underlaying datastore.
func (c *CuratedEpisodeWorker) ProcessDatum(ctx context.Context, datum *messagebusiface.MessageBusMessage) error {
	// episodes are unique by name + date aired
	// check to see if the episode exists
	//	if it does not: add it
	//	else:
	//		check if any of its details of changed
	//		if so, update the details
	//	etc...
	var episodeInfo contracts.EpisodeInfo
	err := json.Unmarshal(datum.Body, &episodeInfo)
	if err != nil {
		nackErr := datum.Acknowledger.Nack(false)
		if nackErr != nil {
			return fmt.Errorf("%v\n%v", err, nackErr)
		}
		return err
	}

	fmt.Println(string(datum.Body))

	return datum.Acknowledger.Ack()
}
