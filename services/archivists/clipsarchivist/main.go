package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/archivists/clipsarchivist"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

func main() {
	log.Println("Connecting to message bus...")
	msgbus, err := amqpadapter.Initialize(context.Background(), "curated_clips", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting clips archivist...")
	clipsArchivist := clipsarchivist.StartWork(context.Background(), msgbus, new(fakeDataStore))

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

type fakeDataStore struct{}

func (f *fakeDataStore) UpsertEpisodeInfo(info contracts.EpisodeInfo) error {
	log.Println(info)
	return nil
}

func (f *fakeDataStore) UpsertClipInfo(clipInfo contracts.ClipInfo) error {
	log.Println(clipInfo)
	return nil
}
