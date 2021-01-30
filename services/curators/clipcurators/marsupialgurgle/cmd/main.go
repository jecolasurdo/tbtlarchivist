package main

import (
	"log"

	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/clipcurators"
	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/curatorbase"
	"github.com/jecolasurdo/tbtlarchivist/pkg/messagebus"
)

func main() {
	msgbus, err := messagebus.Initialize("curated_clips")
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
