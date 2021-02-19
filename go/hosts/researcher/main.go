package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/messagebus/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/engines/researcher"
)

func main() {
	log.Println("Connecting to pending-research queue...")
	pendingQueue, err := amqpadapter.Initialize(context.Background(), "pending_research", amqpadapter.DirectionReceiveOnly)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to completed-research queue...")
	completedQueue, err := amqpadapter.Initialize(context.Background(), "completed_research", amqpadapter.DirectionSendOnly)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting the Research Agent...")
	researchAgent := researcher.StartResearchAgent(context.Background(), pendingQueue, completedQueue, nil)

	log.Println("Running...")
	for {
		select {
		case err := <-researchAgent.Errors:
			log.Println(err)
		case <-researchAgent.Done:
			log.Println("Done")
			return
		}
	}
}
