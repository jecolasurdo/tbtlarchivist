package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/chromedp/chromedp"
	"github.com/jecolasurdo/tbtlarchivist/services/curators/internal/cdp"
	"github.com/jecolasurdo/tbtlarchivist/services/curators/internal/utils"
)

const (
	rawMP3LinkRegex = `href="((?:/[\w+|-]+)+\.mp3)"`
)

var (
	rawMP3LinkRe = regexp.MustCompile(rawMP3LinkRegex)
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
	log.Println("Starting Chrome...")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var rawPageCount string
	err := chromedp.Run(ctx,
		cdp.Logf("Navigating to global search results page..."),
		chromedp.Navigate(`https://www.marsupialgurgle.com/page/1/?s`),
		cdp.Logf("Getting page count..."),
		chromedp.Text(".page-numbers", &rawPageCount, chromedp.ByQuery),
	)
	utils.LogFatalIfErr(err)

	pageCount, err := strconv.Atoi(rawPageCount)
	utils.LogFatalIfErr(err)

	for pageNumber := 1; pageNumber <= pageCount; pageNumber++ {
		var pageBody string
		err := chromedp.Run(ctx,
			chromedp.Navigate("https://www.marsupialgurgle.com/page/%v/?s", pageNumber),
			chromedp.InnerHTML("body", &pageBody, chromedp.BySearch, chromedp.NodeReady),
		)
		utils.LogFatalIfErr(err)

		rawMP3Matches := rawMP3LinkRe.FindAllStringSubmatch(pageBody, -1)
		for i := 0; i < len(rawMP3Matches); i++ {
			// We skip the first "full match" as we're only interested in the
			// submatches.  There should only be one submatch per full match,
			// but we will iterate to be safe.
			for j := 1; j < len(rawMP3Matches[i]); j++ {
				fmt.Println(rawMP3Matches[i][j])
			}
		}

	}
}
