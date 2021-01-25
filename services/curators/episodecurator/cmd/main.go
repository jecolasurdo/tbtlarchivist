package main

import (
	"context"
	"io/ioutil"
	"log"
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

func main() {
	log.Println("Starting Chrome")
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	log.Println("Setting timeout")
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var content string
	var buf []byte
	err := chromedp.Run(ctx,
		logSomething("Navigating to page..."),
		chromedp.Navigate(`https://www.tbtl.net/episodes`),

		// logSomething("Waiting for page to load..."),
		// chromedp.WaitVisible(`body > footer`),

		logSomething("Trying to take a screenshot..."),
		chromedp.CaptureScreenshot(&buf),

		logSomething("Getting content..."),
		chromedp.OuterHTML("#content", &content),
	)
	if err != nil {
		ioutil.WriteFile("screenshot.png", buf, 0o644)
		log.Fatal(err)
	}
	log.Printf("Content:\n%s", content)
}
