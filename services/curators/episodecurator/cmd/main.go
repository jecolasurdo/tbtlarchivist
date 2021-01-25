package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

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
func main() {
	log.Println("Starting Chrome...")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// log.Println("Setting timeout...")
	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	// defer cancel()

	var rawPageCount string
	err := chromedp.Run(ctx,
		logSomething("Navigating to main episodes page..."),
		chromedp.Navigate(`https://www.tbtl.net/episodes`),

		logSomething("Getting page count..."),
		chromedp.Text(".pagination_link-last", &rawPageCount, chromedp.BySearch),
	)
	if err != nil {
		log.Fatal(err)
	}

	pageCount, err := strconv.Atoi(rawPageCount)
	if err != nil {
		log.Fatal(err)
	}

	const hrefRegex = `/episode/\d{4}/\d\d/\d\d/(?:[[:alnum:]]|-)+`
	hrefRe := regexp.MustCompile(hrefRegex)
	episodeLinkList := []string{}
	for pageNumber := 1; pageNumber <= pageCount; pageNumber++ {
		var collectionResults string
		err := chromedp.Run(ctx,
			logSomething(fmt.Sprintf("Scraping page %v...", pageNumber)),
			chromedp.Navigate(fmt.Sprintf("https://www.tbtl.net/episodes/page/%v", pageNumber)),
			chromedp.InnerHTML(".collection_results", &collectionResults, chromedp.NodeVisible, chromedp.BySearch),
		)

		if err != nil {
			log.Fatal(err)
		}

		episodeLinkList = append(episodeLinkList, hrefRe.FindAllString(collectionResults, -1)...)
	}

	const mp3Regex = `/\d{4}/\d\d/\w+\.mp3`
	mp3Re := regexp.MustCompile(mp3Regex)
	for _, episodeLink := range episodeLinkList {
		var nextDataInnerHTML string
		err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("https://www.tbtl.net/%v", episodeLink)),
			chromedp.InnerHTML("#__NEXT_DATA__", &nextDataInnerHTML, chromedp.ByID),
		)

		if err != nil {
			log.Fatal(err)
		}

		mp3Link := mp3Re.FindString(nextDataInnerHTML)
		log.Println("\tmp3 name:", mp3Link)

	}

	fmt.Println("Done.")
}
