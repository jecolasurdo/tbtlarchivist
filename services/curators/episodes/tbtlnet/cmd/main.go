package main

import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/curatorbase"
	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/episodecurators"
	"github.com/jecolasurdo/tbtlarchivist/pkg/messagebus"
)

func main() {
	msgbus, err := messagebus.Initialize(context.Background(), "curated_episodes", 5)
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
