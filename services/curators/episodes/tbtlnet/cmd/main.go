package main

import (
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/curatorbase"
	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/episodecurators"
	"github.com/jecolasurdo/tbtlarchivist/pkg/messagebus"
)

func main() {
	msgbus, err := messagebus.Initialize("curated_episodes")
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
