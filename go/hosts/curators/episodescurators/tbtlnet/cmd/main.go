package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/messagebus/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/engines/curators/curatorbase"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/engines/curators/episodecurators"
)

func main() {
	msgbus, err := amqpadapter.Initialize(context.Background(), "curated_episodes", amqpadapter.DirectionSendOnly)
	if err != nil {
		log.Fatal(err)
	}
	tbtlnet := new(episodecurators.TBTLNet)
	curatorBase := curatorbase.New(tbtlnet, msgbus)
	err = curatorBase.Run()
	if err != nil {
		log.Fatal(err)
	}
}
