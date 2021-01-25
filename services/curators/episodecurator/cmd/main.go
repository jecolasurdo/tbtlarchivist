package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func logSomething(s string) littleLogger {
	return littleLogger{
		s: s,
	}
}

type littleLogger struct {
	s string
}

func (l littleLogger) Do(ctx context.Context) error {
	log.Println(l.s)
	return nil
}

func main() {
	log.Println("Starting Chrome")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	log.Println("Setting timeout")
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// scraping steps:
	// get tbtl.net/episodes
	// wait for the page to be fully loaded
	// identify the total number of pages via the pagination section of the DOM
	// for each page 1..n
	//	note the "teaser link" for each episode
	// for each teaser link
	//	visit the link
	//	record the following to some sort of data store:
	//		episode number
	//		episode part (if multi-part)
	//		episode date
	//		episode title
	//		media link
	//		media type

	var collectionResults string
	var pageCount string
	err := chromedp.Run(ctx,
		logSomething("Navigating to page..."),
		chromedp.Navigate(`https://www.tbtl.net/episodes`),

		logSomething("Getting collection_results..."),
		chromedp.InnerHTML(".collection_results", &collectionResults, chromedp.NodeVisible, chromedp.BySearch),

		logSomething("Getting page count..."),
		chromedp.Text(".pagination_link-last", &pageCount, chromedp.BySearch),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Collection results:\n%s", collectionResults)
	log.Printf("Page count:\n%s", pageCount)
}
