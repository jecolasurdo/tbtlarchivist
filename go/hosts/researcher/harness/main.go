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
			return "/Users/Joe/Documents/code/tbtlarchivist/rust/target/release/analyzerd", nil
		},
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	ctx := context.Background()

	pendingResearchItem := &contracts.PendingResearchItem{
		LeaseId: "test_lease",
		Episode: &contracts.EpisodeInfo{
			MediaUri:  "https://play.publicradio.org/web/o/infinite_guest/tbtl/2021/04/tbtl_20210409_3398_64.mp3",
			MediaType: "mp3",
		},
		Clips: []*contracts.ClipInfo{
			&contracts.ClipInfo{
				MediaUri:  "https://audio.marsupialgurgle.com/audio/andrewandcatchoneinthemiddle-3398.mp3",
				MediaType: "",
			},
		},
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
