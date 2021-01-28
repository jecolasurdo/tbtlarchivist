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
	"github.com/jecolasurdo/tbtlarchivist/pacer"
	"github.com/jecolasurdo/tbtlarchivist/services/curators/internal/utils"
)

const (
	pacingAverage = 4000
	pacingSigma   = 2000

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

func main() {

	log.Printf("Navigating to global search results page...")
	doc, err := htmlquery.LoadURL(`https://www.marsupialgurgle.com/page/1/?s`)
	utils.LogFatalIfErr(err)

	log.Printf("Getting page count...")
	rawPageCount := htmlquery.QuerySelector(doc, pageCountXp).FirstChild.Data
	pageCount, err := strconv.Atoi(rawPageCount)
	utils.LogFatalIfErr(err)
	log.Println(pageCount)

	pace := pacer.SetPace(pacingAverage, pacingSigma, time.Millisecond)
	for pageNumber := 1; pageNumber <= pageCount; pageNumber++ {
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
				continue
			}
			distinctMP3URIs[rawMP3Matches[i][1]] = struct{}{}
		}

		decoratedMP3Matches := mp3WithDescriptionRe.FindAllStringSubmatch(searchHTML, -1)
		distinctDecoratedMP3URIs := map[string]struct{}{}
		for i := 0; i < len(decoratedMP3Matches); i++ {
			if len(decoratedMP3Matches[i]) != 3 {
				continue
			}
			_ = decoratedMP3Matches[i][1] // description to be used later
			mp3URI := decoratedMP3Matches[i][2]
			distinctDecoratedMP3URIs[mp3URI] = struct{}{}
		}

		mp3URIs := []string{}
		// The regex's we're using guarantee that len(distinctMP3URIs) >=
		// len(distinctDecoratedMP3URIs)
		if len(distinctMP3URIs) > len(distinctDecoratedMP3URIs) {
			for m := range distinctMP3URIs {
				if _, found := distinctDecoratedMP3URIs[m]; !found {
					mp3URIs = append(mp3URIs, m)
				}
			}
		}

		fmt.Println("\tundecorated mp3 URIs:", len(mp3URIs))
		fmt.Println("\tdecorated mp3 URIs:", len(distinctDecoratedMP3URIs))

		pace.Wait()
	}
}
