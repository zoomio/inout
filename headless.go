package inout

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func waitForDomElement(ctx context.Context, selector, addr string, verbose bool) (string, error) {
	// create context
	childCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// navigate
	if err := chromedp.Run(childCtx, chromedp.Navigate(addr)); err != nil {
		return "", fmt.Errorf("could not navigate to %s: %v", addr, err)
	}

	var nodes []*cdp.Node
	err := chromedp.Run(childCtx, chromedp.Nodes(selector, &nodes, chromedp.ByQueryAll, chromedp.NodeVisible))
	if err != nil {
		return "", fmt.Errorf("could not get nodes %s: %v", selector, err)
	}

	var sb strings.Builder
	for _, n := range nodes {
		var cnt string
		if err := chromedp.Run(childCtx, chromedp.OuterHTML(n.FullXPath(), &cnt)); err == nil {
			sb.WriteString(cnt)
		}
	}
	return strings.TrimSpace(sb.String()), nil
}
