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
	msgbus, err := amqpadapter.Initialize(context.Background(), "curated_clips", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting clips archivist...")
	clipsArchivist := archivists.StartClipsArchivist(context.Background(), msgbus, new(datastore.FakeDataStorer))

	log.Println("Running...")
	for {
		select {
		case err := <-clipsArchivist.Errors:
			log.Println(err)
		case <-clipsArchivist.Done:
			log.Println("Done")
			return
		}
	}
}
