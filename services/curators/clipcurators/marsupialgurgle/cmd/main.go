package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/jecolasurdo/tbtlarchivist/services/curators/internal/cdp"
)

// clip curation process for marsupial gurgle
// visit page one for a global search
// find the max page count
// for each page
//   extract media uri, originating uri, description
//    - Get a distinct list of all mp3s on page
//    - Build contextualized lists of mp3s
//    - Remove any duplicates (same file name more than once)
func main() {
	log.Println("Starting Chrome...")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	err := chromedp.Run(ctx,
		cdp.Logf("Navigating to main episodes page..."),
		chromedp.Navigate(`https://www.marsupialgurgle.com/page/1/?s`),
	)

	if err != nil {
		log.Fatal(err)
	}
}
