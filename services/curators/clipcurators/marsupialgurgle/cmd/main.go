package main

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
	"github.com/jecolasurdo/tbtlarchivist/services/curators/internal/utils"
)

const (
	marsupialgurgleBaseURI = `https://marsupialgurgle.com`

	pageCountXPath = `//*[@id="bottom-nav-pagination"]/a[3]`

	rawMP3LinkRegex         = `href="((?:/[\w+|-]+)+\.mp3)"`
	mp3WithDescriptionRegex = `(?sU)<p><\w+>(.*)</\w+></p>.*href="((?:/[\w+|-]+)+\.mp3)"`
)

var (
	pageCountXp = xpath.MustCompile(pageCountXPath)

	rawMP3LinkRe         = regexp.MustCompile(rawMP3LinkRegex)
	mp3WithDescriptionRe = regexp.MustCompile(mp3WithDescriptionRegex)
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

	log.Printf("Navigating to global search results page...")
	doc, err := htmlquery.LoadURL(`https://www.marsupialgurgle.com/page/1/?s`)
	utils.LogFatalIfErr(err)

	log.Printf("Getting page count...")
	rawPageCount := htmlquery.QuerySelector(doc, pageCountXp).FirstChild.Data
	pageCount, err := strconv.Atoi(rawPageCount)
	utils.LogFatalIfErr(err)
	log.Println(pageCount)

	fmt.Println("Temporarily starting after page 1 for testing")
	for pageNumber := 113; pageNumber <= pageCount; pageNumber++ {
		paceTime := time.Now().Add(5 * time.Second)

		log.Printf("Scraping page %v of %v...", pageNumber, pageCount)
		resp, err := http.Get(fmt.Sprintf("https://www.marsupialgurgle.com/page/%v/?s", pageNumber))
		utils.LogFatalIfErr(err)
		if resp.StatusCode != 200 {
			log.Fatal(resp.Status)
		}

		searchHTMLBytes, err := ioutil.ReadAll(resp.Body)
		utils.LogFatalIfErr(err)
		searchHTML := string(searchHTMLBytes)

		distinctMP3URIs := map[string]struct{}{}
		rawMP3Matches := rawMP3LinkRe.FindAllStringSubmatch(searchHTML, -1)
		for i := 0; i < len(rawMP3Matches); i++ {
			if len(rawMP3Matches[i]) != 2 {
				log.Println(rawMP3Matches[i])
				continue
			}
			distinctMP3URIs[rawMP3Matches[i][1]] = struct{}{}
		}
		log.Printf("\tDistinct raw mp3 links: %v", len(distinctMP3URIs))

		decoratedMP3Matches := mp3WithDescriptionRe.FindAllStringSubmatch(searchHTML, -1)
		distinctDecoratedMP3URIs := map[string]struct{}{}
		for i := 0; i < len(decoratedMP3Matches); i++ {
			if len(decoratedMP3Matches[i]) != 3 {
				log.Println(decoratedMP3Matches[i])
				continue
			}
			_ = decoratedMP3Matches[i][1] // description to be used later
			mp3URI := decoratedMP3Matches[i][2]
			distinctDecoratedMP3URIs[mp3URI] = struct{}{}
		}
		log.Printf("\tDistinct decorated mp3s links: %v", len(distinctDecoratedMP3URIs))

		if len(distinctDecoratedMP3URIs) != len(distinctMP3URIs) {

			// pages 83 and 112 contain a variant in the html formatting which is
			// <p><strong><br><a href...
			// whereas the following is more typical for the site
			// <p><strong> ... </p></strong><p><a href...
			// Page 112 example /audio/lukeandrewdoyouneedsomealcohol-2748.mp3
			log.Printf("Mismatch on page %v", pageNumber)

			for m := range distinctMP3URIs {
				if _, found := distinctDecoratedMP3URIs[m]; !found {
					fmt.Println(m)
				}
			}
		}

		now := time.Now()
		if now.Before(paceTime) {
			waitDuration := paceTime.Sub(now)
			log.Printf("Pacing (%v)...", waitDuration)
			time.Sleep(waitDuration)
		}
	}
}
