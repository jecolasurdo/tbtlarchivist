package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/researcher/agent"
)

func main() {
	log.Println("Connecting to pending-research queue...")
	pendingQueue, err := amqpadapter.Initialize(context.Background(), "pending_research", 1)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to completed-research queue...")
	completedQueue, err := amqpadapter.Initialize(context.Background(), "completed_research", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting the Research Agent...")
	researchAgent := agent.StartResearchAgent(context.Background(), pendingQueue, completedQueue)

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
