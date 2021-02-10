package main

import (
	"context"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore/adapters/mariadbadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/archivists"
)

func main() {
	log.Println("Connecting to database...")
	dbconfig := mysql.NewConfig()
	dbconfig.Addr = "127.0.0.1:3306"
	dbconfig.DBName = "tbtlarchivist"
	dbconfig.User = "root"
	log.Println(dbconfig.FormatDSN())
	db, err := mariadbadapter.New(dbconfig.FormatDSN(), 60*time.Second, 5, 5).Connect()
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
		case err := <-clipsArchivist.Errors:
			log.Println(err)
		case <-clipsArchivist.Done:
			log.Println("Done")
			return
		}
	}
}
