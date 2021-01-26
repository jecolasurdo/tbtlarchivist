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

	// NumberIsApplicable describes whether or not an episode number applies to
	// this particular episode. For instance, some episodes, such as the "no
	// point conversions" are non-canonical, in that they are not included in
	// the episode tally. This field is true unless an episode is explicitly
	// known to be non-canonical.
	CanonicalNumberIsApplicable bool

	// CanonicalNumberDerivation describes how the episode number was derived.
	// Episode numbers can be inferred from various means such as being
	// extracted from the episode title, extracted from the mp3 file name, etc.
	CanonicalNumberDerivation string

	// CanonicalNumber is the episode number for any "canonical" episodes.
	// "canonical episodes" are any that are not otherwise excluded by some
	// rule such as "no point conversion" episodes, which are non-canonical.
	// This value will be greater than 0 for all canonical episodes that have
	// had their number infered.  This value will be -1 for any canonical
	// episodes for whom a number could not be inferred.  This value will be -2
	// for any non-canonical episodes.
	CanonicalNumber int

	// Part represents the episde segment. If an episode was uploaded in
	// multiple segments, each segment is numbered in order here.
	Part int

	// Parts represents how many segments exist for an episode. If an episode
	// was uploaded in multiple segments, this value represents the total
	// uploaded.
	Parts int

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
			err := chromedp.Run(ctx,
				chromedp.Navigate(fmt.Sprintf("https://www.tbtl.net/%v", episodeLink)),
				chromedp.InnerHTML("#__NEXT_DATA__", &nextDataInnerHTML, chromedp.ByID),
				chromedp.TextContent(".hdg", &title, chromedp.BySearch),
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

			episodeInfo := EpisodeInfo{
				DateCurated:     time.Now().UTC(),
				DateAired:       time.Now().UTC(),
				CanonicalNumber: -1,
				Part:            -1,
				Duration:        time.Duration(durationMS) * time.Millisecond,
				Title:           title,
				Description:     "",
				MediaURI:        mediaURI,
				MediaType:       mediaType,
			}

			fmt.Println(episodeInfo)
		}
	}

	fmt.Println("Done.")
}
