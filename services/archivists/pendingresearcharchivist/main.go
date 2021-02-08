package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/archivists"
)

func main() {

	log.Println("Connecting to message bus...")
	msgbus, err := amqpadapter.Initialize(context.Background(), "pending_research", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting pending-research archivist...")
	pendingResearchArchivist := archivists.StartPendingResearchArchivist(context.Background(), msgbus, new(datastore.FakeDataStorer))

	log.Println("Running...")
	for {
		select {
		case err := <-pendingResearchArchivist.Errors:
			log.Println(err)
		case <-pendingResearchArchivist.Done:
			log.Println("Done")
			return
		}
	}
}
