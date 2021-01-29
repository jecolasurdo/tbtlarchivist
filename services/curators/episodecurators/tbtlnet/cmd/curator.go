package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/jecolasurdo/tbtlarchivist/pacer"
	"github.com/jecolasurdo/tbtlarchivist/services/curators/contracts"
	"github.com/jecolasurdo/tbtlarchivist/services/curators/internal/cdp"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"curated_episodes", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare a queue")

	episodeInfoSource, errorSource := Curate()
	channelsOpen := true
	for channelsOpen {
		select {
		case episodeInfo := <-episodeInfoSource:
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(episodeInfo.String()),
				})
			failOnError(err, "Failed to publish a message")
		case err, isOpen := <-errorSource:
			channelsOpen = isOpen
			fmt.Println(err)
		}
	}
}

const (
	scraperName = `tbtl.net scraper`

	// Duration is extracted from the __NEXT_DATA__ structure within the DOM.
	// __NEXT_DATA__ may contain data for m4a files in addition to mp3s, and
	// the m4a durations may not match that of the mp3s.  To address this,
	// durationRegex qualifies the duration_ms field as being directly
	// preceeded by the value `mp3\",\"`. This ensures that the duration data
	// is associated with the correct file, but is admittedly a little fragile
	// in that it presumes field order is fixed.
	durationRegex = `mp3\\",\\"duration_ms\\":(\d+)`
	hrefRegex     = `/episode/\d{4}/\d\d/\d\d/(?:[[:alnum:]]|-)+`
	mp3Regex      = `https://(?:(?:\w+|-|\.)+/)+\d{4}/\d{1,2}/(?:\d{1,2}/)?(?:\w+|-)+\.mp3`

	unreplacedUAToken = "unreplaced_ua"
	userAgent         = "web"
)

var durationRe = regexp.MustCompile(durationRegex)
var hrefRe = regexp.MustCompile(hrefRegex)
var mp3Re = regexp.MustCompile(mp3Regex)

// Curate initializes the scraper and returns two channels, one providing a
// stream of episode information that has been scraped, and the other containing
// any errors that have been emited by the process.
func Curate() (<-chan *contracts.EpisodeInfo, <-chan error) {
	episodeInfoSource := make(chan *contracts.EpisodeInfo)
	errorSource := make(chan error)

	go func() {
		defer close(episodeInfoSource)
		defer close(errorSource)

		log.Println("Starting Chrome (headless)...")
		ctx, cancel := chromedp.NewContext(
			context.Background(),
			chromedp.WithLogf(log.Printf),
		)
		defer cancel()

		var rawPageCount string
		err := chromedp.Run(ctx,
			cdp.Logf("Navigating to main episodes page..."),
			chromedp.Navigate(`https://www.tbtl.net/episodes`),

			cdp.Logf("Getting page count..."),
			chromedp.Text(".pagination_link-last", &rawPageCount, chromedp.BySearch),
		)
		if err != nil {
			errorSource <- err
			return
		}

		pageCount, err := strconv.Atoi(rawPageCount)
		if err != nil {
			errorSource <- err
			return
		}

		log.Println("Scraping...")
		pace := pacer.SetPace(1000, 300, time.Millisecond)
		for pageNumber := 1; pageNumber <= pageCount; pageNumber++ {
			var collectionResults string
			err := chromedp.Run(ctx,
				chromedp.Navigate(fmt.Sprintf("https://www.tbtl.net/episodes/page/%v", pageNumber)),
				chromedp.InnerHTML(".collection_results", &collectionResults, chromedp.NodeVisible, chromedp.BySearch),
			)
			if err != nil {
				errorSource <- err
				return
			}

			episodeLinkList := hrefRe.FindAllString(collectionResults, -1)
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
					errorSource <- err
					return
				}

				mediaURI := mp3Re.FindString(nextDataInnerHTML)
				if mediaURI == "" {
					log.Fatal("Media URI Not Identified")
				}
				mediaURI = strings.Replace(mediaURI, unreplacedUAToken, userAgent, -1)
				mediaType := mediaURI[len(mediaURI)-3:]

				rawDuration := durationRe.FindStringSubmatch(nextDataInnerHTML)
				if len(rawDuration) < 2 {
					log.Fatal("Unable to extract duration for episode.")
				}
				durationMS, err := strconv.Atoi(rawDuration[1])
				if err != nil {
					errorSource <- err
					return
				}

				dateAired, err := time.Parse("January 2, 2006", rawDate)
				if err != nil {
					errorSource <- err
					return
				}

				episodeInfoSource <- &contracts.EpisodeInfo{
					CuratorInformation: scraperName,
					DateCurated:        time.Now().UTC(),
					DateAired:          dateAired,
					Duration:           time.Duration(durationMS) * time.Millisecond,
					Title:              title,
					Description:        description,
					MediaURI:           mediaURI,
					MediaType:          mediaType,
				}

				pace.Wait()
			}
			pace.Wait()
		}

		fmt.Println("Done.")
	}()

	return episodeInfoSource, errorSource
}
