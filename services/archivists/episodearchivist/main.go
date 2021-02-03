package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/archivists/episodearchivist"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

func main() {
	log.Println("Connecting to message bus...")
	msgbus, err := amqpadapter.Initialize(context.Background(), "curated_episodes", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting episode archivist...")
	episodeArchivist := episodearchivist.StartWork(context.Background(), msgbus, new(fakeDataStore))

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

type fakeDataStore struct{}

func (f *fakeDataStore) UpsertEpisodeInfo(info contracts.EpisodeInfo) error {
	log.Println(info)
	return nil
}

func (f *fakeDataStore) UpsertClipInfo(clipInfo contracts.ClipInfo) error {
	log.Println(clipInfo)
	return nil
}
