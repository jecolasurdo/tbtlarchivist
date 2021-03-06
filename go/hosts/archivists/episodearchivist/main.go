package main

import (
	"context"
	"log"
	"time"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/datastore/adapters/mariadbadapter"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/messagebus/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/engines/archivists"
)

func main() {
	log.Println("Connecting to database...")
	dbconfig := &mariadbadapter.Config{
		Addr:                  "127.0.0.1:3306",
		DBName:                "tbtlarchivist",
		User:                  "root",
		MaxConnectionLifetime: 60 * time.Second,
		MaxOpenConnections:    5,
		MaxIdleConnections:    5,
	}
	db, err := mariadbadapter.New(dbconfig).Connect()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to message bus...")
	msgbus, err := amqpadapter.Initialize(context.Background(), "curated_episodes", amqpadapter.DirectionSendOnly)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting episode archivist...")
	episodeArchivist := archivists.StartEpisodesArchivist(context.Background(), msgbus, db)

	log.Println("Running...")
	for {
		select {
		case err, open := <-episodeArchivist.Errors:
			if !open {
				break
			}
			log.Println(err)
		case <-episodeArchivist.Done:
			log.Println("Done")
			return
		}
	}
}
