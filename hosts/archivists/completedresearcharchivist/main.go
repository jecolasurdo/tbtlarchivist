package main

// import (
// 	"context"
// 	"log"

// 	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore/adapters/fakedatastore"
// 	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus/adapters/amqpadapter"
// 	"github.com/jecolasurdo/tbtlarchivist/pkg/archivists"
// )

// func main() {

// 	log.Println("Connecting to message bus...")
// 	msgbus, err := amqpadapter.Initialize(context.Background(), "completed_research", amqpadapter.DirectionSendOnly)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Starting completed-research archivist...")
// 	completedResearchArchivist := archivists.StartCompletedResearchArchivist(context.Background(), msgbus, new(fakedatastore.FakeDataStorer))

// 	log.Println("Running...")
// 	for {
// 		select {
// 		case err := <-completedResearchArchivist.Errors:
// 			log.Println(err)
// 		case <-completedResearchArchivist.Done:
// 			log.Println("Done")
// 			return
// 		}
// 	}
// }
