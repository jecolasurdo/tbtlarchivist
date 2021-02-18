package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus/adapters/amqpadapter"
	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/clipcurators"
	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/curatorbase"
)

func main() {
	msgbus, err := amqpadapter.Initialize(context.Background(), "curated_clips", amqpadapter.DirectionSendOnly)
	if err != nil {
		log.Fatal(err)
	}

	marsuplialGurgle := new(clipcurators.MarsupialGurgle)
	curatorBase := curatorbase.New(marsuplialGurgle, msgbus)

	err = curatorBase.Run()
	if err != nil {
		log.Fatal(err)
	}
}
