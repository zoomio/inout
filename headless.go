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

	// get project link text
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(selector, &nodes)); err != nil {
		return "", err
	}

	var b strings.Builder
	for _, n := range nodes {
		fmt.Fprintf(&b, "%s", getNodeText(n, "p"))
	}
	return strings.TrimSpace(b.String()), nil
}

func getNodeText(node *cdp.Node, parent string) string {
	if node.NodeValue != "" {
		return fmt.Sprintf("<%s>%s</%s>", parent, node.NodeValue, parent)
	}

	if node.ChildNodeCount == 0 {
		return ""
	}

	var b strings.Builder
	for _, n := range node.Children {
		text := getNodeText(n, node.LocalName)
		if text != "" {
			fmt.Fprintf(&b, "%s", text)
		}
	}

	return strings.TrimSpace(b.String())
}
