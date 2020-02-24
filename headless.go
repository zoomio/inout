package inout

import (
	"context"

	"github.com/chromedp/chromedp"
)

func waitForDomElement(ctx context.Context, selector, host string, verbose bool) (string, error) {
	// create context
	childCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var content string
	err := chromedp.Run(childCtx, visible(host, selector, &content))
	if err != nil {
		return "", err
	}

	return content, nil
}

func visible(host, selector string, content *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(host),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(selector, content),
	}
}
