package clipcurators

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
	"github.com/jecolasurdo/tbtlarchivist/pkg/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	scraperName = `marsupialgurgle scraper`

	pacingAverage = 4000
	pacingSigma   = 2000

	marsupialgurgleBaseURI = `https://marsupialgurgle.com`

	pageCountXPath = `//*[@id="bottom-nav-pagination"]/a[3]`

	// The scraper assumes the rawMP3LinkRegex is a subset of the
	// mp3WithDescriptionRegex. This ensures that the number of rawMP3 links
	// found must be greater or equal to the number of links found with
	// descriptions.
	rawMP3LinkRegex         = `href="((?:/[\w+|-]+)+\.mp3)"`
	mp3WithDescriptionRegex = `(?sU)<p><\w+>(.*)</\w+></p>.*href="((?:/[\w+|-]+)+\.mp3)"`
)

var (
	pageCountXp = xpath.MustCompile(pageCountXPath)

	rawMP3LinkRe         = regexp.MustCompile(rawMP3LinkRegex)
	mp3WithDescriptionRe = regexp.MustCompile(mp3WithDescriptionRegex)
)

// MarsupialGurgle is a curator that extracts data from
// www.marsupialgurgle.net.
type MarsupialGurgle struct {
}

// Curate initializes the scraper and returns two channels, one providing a
// stream of clip information that has been scraped, and the other containing
// any errors that have been emited by the process.
func (m *MarsupialGurgle) Curate() (<-chan protoreflect.ProtoMessage, <-chan error) {
	clipInfoSource := make(chan protoreflect.ProtoMessage)
	errorSource := make(chan error)

	go func() {
		defer close(clipInfoSource)
		defer close(errorSource)

		log.Printf("Navigating to global search results page...")
		uri := `https://www.marsupialgurgle.com/page/1/?s`
		doc, err := htmlquery.LoadURL(uri)
		if err != nil {
			errorSource <- fmt.Errorf("error loading %v (%v)", uri, err)
			return
		}

		log.Printf("Getting page count...")
		rawPageCount := htmlquery.QuerySelector(doc, pageCountXp).FirstChild.Data
		pageCount, err := strconv.Atoi(rawPageCount)
		if err != nil {
			errorSource <- fmt.Errorf("error while extracting page count: %v", err)
			return
		}

		pace := utils.SetNormalPace(pacingAverage, pacingSigma, time.Millisecond)

		// We visit the pages in random order to increase the breadth of each
		// search, in case the search gets terminated before all pages have
		// been visited.
		shuffledPages := utils.GetShuffledIntList(pageCount)
		for _, pageNumber := range shuffledPages {
			log.Printf("Scraping page %v of %v...", pageNumber, pageCount)
			pageURI := fmt.Sprintf("https://www.marsupialgurgle.com/page/%v/?s", pageNumber)
			resp, err := http.Get(pageURI)
			if err != nil {
				errorSource <- fmt.Errorf("error accessing page %v (%v)", pageURI, err)
				continue
			}
			if resp.StatusCode != 200 {
				errorSource <- fmt.Errorf("received non-200 response when requesting page %v", pageURI)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errorSource <- fmt.Errorf("error occured when reading the response body for page %v (%v)", pageURI, err)
				continue
			}

			mp3s := extractMP3s(string(body))
			decoratedMP3s := extractDecoratedMP3s(string(body))

			// Since the rawMP3Link regex is a subset of the mp3WithDescription
			// regex, we are guaranteed that len(distinctMP3URIs) >=
			// len(distinctDecoratedMP3URIs)
			if len(mp3s) > len(decoratedMP3s) {
				for mp3 := range decoratedMP3s {
					if _, found := mp3s[mp3]; !found {
						delete(mp3s, mp3)
					}
				}
			}

			for _, clipInfo := range mp3s {
				clipInfoSource <- clipInfo
			}

			for _, clipInfo := range decoratedMP3s {
				clipInfoSource <- clipInfo
			}

			pace.Wait()
		}
	}()

	return clipInfoSource, errorSource
}

func extractMP3s(body string) map[string]*contracts.ClipInfo {
	distinctMP3URIs := map[string]*contracts.ClipInfo{}
	rawMP3Matches := rawMP3LinkRe.FindAllStringSubmatch(body, -1)
	for i := 0; i < len(rawMP3Matches); i++ {
		if len(rawMP3Matches[i]) != 2 {
			continue
		}
		mp3URI := rawMP3Matches[i][1]
		now := timestamppb.New(time.Now().UTC())
		distinctMP3URIs[mp3URI] = &contracts.ClipInfo{
			InitialDateCurated: now,
			LastDateCurated:    now,
			CuratorInformation: scraperName,
			Title:              mp3URI,
			Description:        "",
			MediaUri:           mp3URI,
			MediaType:          "mp3",
			Priority:           0,
		}
	}
	return distinctMP3URIs
}

func extractDecoratedMP3s(body string) map[string]*contracts.ClipInfo {
	decoratedMP3Matches := mp3WithDescriptionRe.FindAllStringSubmatch(body, -1)
	distinctDecoratedMP3URIs := map[string]*contracts.ClipInfo{}
	for i := 0; i < len(decoratedMP3Matches); i++ {
		if len(decoratedMP3Matches[i]) != 3 {
			continue
		}
		mp3URI := decoratedMP3Matches[i][2]
		now := timestamppb.New(time.Now().UTC())
		distinctDecoratedMP3URIs[mp3URI] = &contracts.ClipInfo{
			InitialDateCurated: now,
			LastDateCurated:    now,
			CuratorInformation: scraperName,
			Title:              mp3URI,
			Description:        decoratedMP3Matches[i][1],
			MediaUri:           mp3URI,
			MediaType:          "mp3",
			Priority:           0,
		}
	}
	return distinctDecoratedMP3URIs
}
