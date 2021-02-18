package main

import (
	"context"
	"log"
	"time"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore/adapters/mariadbadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/archivists"
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
	msgbus, err := amqpadapter.Initialize(context.Background(), "pending_research", amqpadapter.DirectionSendOnly)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting pending-research archivist...")
	pendingResearchArchivist := archivists.StartPendingResearchArchivist(context.Background(), msgbus, db)

	log.Println("Running...")
	for {
		select {
		case err, open := <-pendingResearchArchivist.Errors:
			if !open {
				break
			}
			log.Println(err)
		case <-pendingResearchArchivist.Done:
			log.Println("Done")
			return
		}
	}
}
