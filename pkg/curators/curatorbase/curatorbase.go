package curatorbase

import (
	"encoding/json"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"
	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/curatoriface"
)

// CuratorBase manages communication between a curator and a message bus.
type CuratorBase struct {
	curator    curatoriface.Curator
	messageBus messagebus.Sender
}

// New returns a reference to a new CuratorBase instance.
func New(curator curatoriface.Curator, messageBus messagebus.Sender) *CuratorBase {
	if curator == nil {
		panic("curator must not be nil")
	}

	if messageBus == nil {
		panic("messageBus must not be nil")
	}

	return &CuratorBase{
		curator:    curator,
		messageBus: messageBus,
	}
}

// Run calls the underlaying curator's Curate method, and begins sending the
// resulting data stream to the message bus.
func (c *CuratorBase) Run() (err error) {
	resultSource, errorSource := c.curator.Curate()
poll:
	for {
		select {
		case result := <-resultSource:
			jsonBytes, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				break poll
			}
			err = c.messageBus.Send(jsonBytes)
			if err != nil {
				break poll
			}
		case err = <-errorSource:
			break poll
		}
	}

	return err
}
