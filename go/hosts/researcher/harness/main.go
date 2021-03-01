package main

// This "harness" is used for manually experimenting with interop between a
// downstream analyst cli and an Adapter. The normal host doesn't call the
// Adapter directly, as the Adapter is typically called via an Analyzer
// accessor.
import (
	"context"
	"log"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/analyst/adapters/rustanalyst"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
)

func main() {
	adapter := rustanalyst.Adapter{
		PathResolver: func() (string, error) {
			return "/Users/Joe/Documents/code/tbtlarchivist/rust/analyst/target/release/cli", nil
		},
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	ctx := context.Background()

	pendingResearchItem := &contracts.PendingResearchItem{
		LeaseId: "test",
		Episode: &contracts.EpisodeInfo{
			MediaUri:  "",
			MediaType: "mp3",
		},
		Clips: []*contracts.ClipInfo{},
	}

	adapter.Run(ctx, pendingResearchItem)

	for {
		select {
		case <-adapter.Done():
			log.Println("Done")
			return
		case err := <-adapter.Errors():
			if err != nil {
				log.Println(err)
			}
		case work := <-adapter.CompletedWorkItems():
			if work != nil {
				log.Println(work)
			}
		}
	}
}
