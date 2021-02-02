package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/curatorbase"
	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/episodecurators"
)

func main() {
	msgbus, err := amqpadapter.Initialize(context.Background(), "curated_episodes", 5)
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