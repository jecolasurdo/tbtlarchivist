package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore"
	"github.com/jecolasurdo/tbtlarchivist/pkg/archivists"
)

func main() {
	log.Println("Connecting to message bus...")
	msgbus, err := amqpadapter.Initialize(context.Background(), "curated_episodes", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting episode archivist...")
	episodeArchivist := archivists.StartEpisodesArchivist(context.Background(), msgbus, new(datastore.FakeDataStorer))

	log.Println("Running...")
	for {
		select {
		case err := <-episodeArchivist.Errors:
			log.Println(err)
		case <-episodeArchivist.Done:
			log.Println("Done")
			return
		}
	}
}
