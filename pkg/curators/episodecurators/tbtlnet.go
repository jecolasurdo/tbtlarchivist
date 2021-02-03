package episodecurators

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
	"github.com/jecolasurdo/tbtlarchivist/pkg/curators/utils"
)

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

// TBTLNet is a curator that extracts episode data from www.tbtl.net.
type TBTLNet struct{}

// Curate initializes the scraper and returns two channels, one providing a
// stream of episode information that has been scraped, and the other containing
// any errors that have been emited by the process.
func (t *TBTLNet) Curate() (<-chan interface{}, <-chan error) {
	episodeInfoSource := make(chan interface{})
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
			utils.Logf("Navigating to main episodes page..."),
			chromedp.Navigate(`https://www.tbtl.net/episodes`),

			utils.Logf("Getting page count..."),
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
		pace := utils.SetPace(1000, 300, time.Millisecond)

		// We visit the pages in random order to increase the breadth of each
		// search, in case the search gets terminated before all pages have
		// been visited.
		shuffledPages := utils.GetShuffledIntList(pageCount)
		for _, pageNumber := range shuffledPages {
			var collectionResults string
			uri := fmt.Sprintf("https://www.tbtl.net/episodes/page/%v", pageNumber)
			err := chromedp.Run(ctx,
				chromedp.Navigate(uri),
				chromedp.InnerHTML(".collection_results", &collectionResults, chromedp.NodeVisible, chromedp.BySearch),
			)
			if err != nil {
				errorSource <- fmt.Errorf("error while accessing %v (%v)", uri, err)
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
					errorSource <- fmt.Errorf("error while extracting episode page data. %v", err)
					continue
				}

				mediaURI := mp3Re.FindString(nextDataInnerHTML)
				if mediaURI == "" {
					errorSource <- fmt.Errorf("unable to extract media URI. %v", episodeLink)
					continue
				}
				mediaURI = strings.Replace(mediaURI, unreplacedUAToken, userAgent, -1)
				mediaType := mediaURI[len(mediaURI)-3:]

				rawDuration := durationRe.FindStringSubmatch(nextDataInnerHTML)
				if len(rawDuration) < 2 {
					errorSource <- fmt.Errorf("Unable to extract duration for episode. %v", episodeLink)
					continue
				}
				durationMS, err := strconv.Atoi(rawDuration[1])
				if err != nil {
					errorSource <- fmt.Errorf("Unable to parse episode duration. (%v) %v", err, episodeLink)
					continue
				}

				dateAired, err := time.Parse("January 2, 2006", rawDate)
				if err != nil {
					errorSource <- fmt.Errorf("Unable to parse date aired. (%v) %v", err, episodeLink)
					continue
				}

				episodeInfoSource <- contracts.EpisodeInfo{
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
