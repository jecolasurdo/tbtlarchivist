package utils

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

// Logf is an adapter for `log.Printf` which can be used as a chromedp.Action.
func Logf(format string, v ...interface{}) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Printf(format, v...)
		return nil
	})
}
