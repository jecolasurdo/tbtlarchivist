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
	// DateCurated represents the date that the curator service found and
	// analyzed the episode.
	DateCurated time.Time

	// CuratorInformation provides information about the utility that extracted
	// this information.
	CuratorInformation string

	// DateAired is the date that the episode was originally aired.
	DateAired time.Time

	// Duration is the length of the episode.
	Duration time.Duration

	// Title is the name of the episode.
	Title string

	// Description is the episode description.
	Description string

	// MediaURI is a URI for where the episode media can be accessed.
	MediaURI string

	// MediaType is the media type for the episode (such as mp3, etc).
	MediaType string
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

		// The __NEXT_DATA__ structure may contain data for m4a files in
		// addition to mp3s, and the m4a durations may not match that of the
		// mp3s.  To address this, durationRegex qualifies the duration_ms
		// field as being directly preceeded by the value `mp3\",\"`. This
		// ensures that the duration data is associated with the correct file,
		// but is admittedly a little fragile in that it presumes field order
		// is fixed.
		const durationRegex = `mp3\\",\\"duration_ms\\":(\d+)`
		durationRe := regexp.MustCompile(durationRegex)

		const mp3Regex = `https://(?:(?:\w+|-|\.)+/)+\d{4}/\d{1,2}/(?:\d{1,2}/)?(?:\w+|-)+\.mp3`
		mp3Re := regexp.MustCompile(mp3Regex)
		for _, episodeLink := range episodeLinkList {
			var nextDataInnerHTML string
			var title string
			var description string
			var rawDate string
			err := chromedp.Run(ctx,
				chromedp.Navigate(fmt.Sprintf("https://www.tbtl.net/%v", episodeLink)),
				chromedp.InnerHTML("#__NEXT_DATA__", &nextDataInnerHTML, chromedp.ByID),
				chromedp.TextContent(".hdg", &title, chromedp.BySearch),
				chromedp.TextContent("body > div > main > div > section > div > article > div > div > div > p", &description, chromedp.ByQuery),
				chromedp.TextContent(".content_date", &rawDate, chromedp.BySearch),
			)

			if err != nil {
				log.Fatal(err)
			}

			mediaURI := mp3Re.FindString(nextDataInnerHTML)

			if mediaURI == "" {
				log.Fatal("Media URI Not Identified")
			}

			const unreplacedUAToken = "unreplaced_ua"
			const userAgent = "web"
			mediaURI = strings.Replace(mediaURI, unreplacedUAToken, userAgent, -1)

			mediaType := mediaURI[len(mediaURI)-3:]

			rawDuration := durationRe.FindStringSubmatch(nextDataInnerHTML)

			if len(rawDuration) < 2 {
				log.Fatal("Unable to extract duration for episode.")
			}

			durationMS, err := strconv.Atoi(rawDuration[1])
			if err != nil {
				log.Fatal(err)
			}

			dateAired, err := time.Parse("January 2, 2006", rawDate)
			if err != nil {
				log.Fatal(err)
			}

			episodeInfo := EpisodeInfo{
				CuratorInformation: "tbtl.net scraper",
				DateCurated:        time.Now().UTC(),
				DateAired:          dateAired,
				Duration:           time.Duration(durationMS) * time.Millisecond,
				Title:              title,
				Description:        description,
				MediaURI:           mediaURI,
				MediaType:          mediaType,
			}

			fmt.Println(episodeInfo.String() + ",")
		}
	}

	fmt.Println("Done.")
}
