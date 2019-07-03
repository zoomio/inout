package inout

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func waitForDomElement(selector, host string) (string, error) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// navigate
	if err := chromedp.Run(ctx, chromedp.Navigate(host)); err != nil {
		return "", err
	}

	// wait visible
	if err := chromedp.Run(ctx, chromedp.WaitVisible(selector)); err != nil {
		return "", err
	}

	// get text
	// var res string
	// if err := chromedp.Run(ctx, chromedp.Text(selector, &res, chromedp.NodeVisible)); err != nil {
	// 	return "", err
	// }

	// get project link text
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(selector, &nodes)); err != nil {
		return "", err
	}

	var b strings.Builder
	for _, n := range nodes {
		fmt.Printf("- %s\n", n.NodeValue)
		fmt.Fprintf(&b, "%s ", n.NodeValue)
	}
	return strings.TrimSpace(b.String()), nil

	// return strings.TrimSpace(res), nil
}
