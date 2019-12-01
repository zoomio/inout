package inout

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func waitForDomElement(ctx context.Context, selector, host string, verbose bool) (string, error) {
	// create context
	childCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if verbose {
		fmt.Printf("navigating to %s\n", host)
	}

	// navigate
	if err := chromedp.Run(childCtx, chromedp.Navigate(host)); err != nil {
		return "", err
	}

	if verbose {
		fmt.Printf("waiting for \"%s\" to be visible\n", selector)
	}

	// wait to be visible
	if err := chromedp.Run(childCtx, chromedp.WaitVisible(selector)); err != nil {
		return "", err
	}

	if verbose {
		fmt.Println("collecting nodes")
	}

	// collect relevant nodes
	var nodes []*cdp.Node
	if err := chromedp.Run(childCtx, chromedp.Nodes(selector, &nodes)); err != nil {
		return "", err
	}

	if verbose {
		fmt.Println("collecting text from nodes")
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
