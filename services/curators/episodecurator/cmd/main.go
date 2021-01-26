package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
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

// EpisodeInfo contains information about an episode.
type EpisodeInfo struct {
	DateFound   time.Time
	DateAired   time.Time
	Number      int
	Part        int
	Duration    time.Duration
	Title       string
	Description string
	MediaLink   string
	MediaType   string
}

func (e EpisodeInfo) String() string {
	jsonBytes, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

func main() {
	log.Println("Starting Chrome...")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

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
	for pageNumber := 1; pageNumber <= pageCount; pageNumber++ {
		fmt.Printf("Page: %v\n", pageNumber)

		var collectionResults string
		err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("https://www.tbtl.net/episodes/page/%v", pageNumber)),
			chromedp.InnerHTML(".collection_results", &collectionResults, chromedp.NodeVisible, chromedp.BySearch),
		)

		if err != nil {
			log.Fatal(err)
		}

		episodeLinkList := hrefRe.FindAllString(collectionResults, -1)

		const durationRegex = `mp3\\",\\"duration_ms\\":(\d+)`
		durationRe := regexp.MustCompile(durationRegex)

		const mp3Regex = `https://(?:(?:\w+|-|\.)+/)+\d{4}/\d{1,2}/(?:\d{1,2}/)?(?:\w+|-)+\.mp3`
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

			mediaLink := mp3Re.FindString(nextDataInnerHTML)

			if mediaLink == "" {
				log.Fatal("Media Link Not Identified")
			}

			const unreplacedUAToken = "unreplaced_ua"
			const userAgent = "web"
			mediaLink = strings.Replace(mediaLink, unreplacedUAToken, userAgent, -1)

			mediaType := mediaLink[len(mediaLink)-3:]

			rawDuration := durationRe.FindStringSubmatch(nextDataInnerHTML)
			durationMS, err := strconv.Atoi(rawDuration[1])
			if err != nil {
				log.Fatal(err)
			}

			episodeInfo := EpisodeInfo{
				DateFound:   time.Now().UTC(),
				DateAired:   time.Now().UTC(),
				Number:      -1,
				Part:        -1,
				Duration:    time.Duration(durationMS) * time.Millisecond,
				Title:       "",
				Description: "",
				MediaLink:   mediaLink,
				MediaType:   mediaType,
			}

			fmt.Println(episodeInfo)
		}
	}

	fmt.Println("Done.")
}
