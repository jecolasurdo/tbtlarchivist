package main

import (
	"context"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	msgbus, err := amqpadapter.Initialize(context.Background(), "curated_clips", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting clips archivist...")
	clipsArchivist := archivists.StartClipsArchivist(context.Background(), msgbus, db)

	log.Println("Running...")
	for {
		select {
		case err, open := <-clipsArchivist.Errors:
			if !open {
				break
			}
			log.Println(err)
		case <-clipsArchivist.Done:
			log.Println("Done")
			return
		}
	}
}
