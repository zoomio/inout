package inout

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

type headlesResult struct {
	htmlDoc  string
	imgBytes []byte
}

func headless(ctx context.Context, c *config) (*headlesResult, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
	)
	if c.userAgent != "" {
		opts = append(opts, chromedp.UserAgent(c.userAgent))
	}
	// Create an allocator
	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocatorCancel()

	// Create a new context with the allocator
	childCtx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	var res strings.Builder
	var img []byte
	if err := chromedp.Run(
		childCtx,
		// chromedp.Emulate(device.IPhone7), uncomment to emulate devices
		chromeTasks(c, &res, 90, &img)); err != nil {
		err = fmt.Errorf("error in running headless to %s: %w", c.source, err)
		return nil, err
	}

	return &headlesResult{htmlDoc: strings.TrimSpace(res.String()), imgBytes: img}, nil
}

// chromeTasks ...
func chromeTasks(c *config, res *strings.Builder, quality int, buf *[]byte) chromedp.Tasks {
	tasks := []chromedp.Action{chromedp.ActionFunc(func(ctx context.Context) error {
		c := chromedp.FromContext(ctx)
		_, err := target.CreateBrowserContext().Do(cdp.WithExecutor(ctx, c.Browser))
		return err
	})}

	if c.screenshot {
		tasks = append(tasks, chromedp.EmulateViewport(1920, 2000))
	}
	tasks = append(tasks, chromedp.Navigate(c.source))
	if c.waitUntil > 0 {
		tasks = append(tasks, chromedp.Sleep(c.waitUntil))
	}
	if len(c.waitFor) > 0 {
		tasks = append(tasks, chromedp.WaitReady(c.waitFor))
	}
	if len(c.query) > 0 {
		var nodes []*cdp.Node
		tasks = append(tasks,
			chromedp.WaitReady(c.query),
			chromedp.Nodes(c.query, &nodes),
			chromedp.ActionFunc(func(c context.Context) error {
				for _, v := range nodes {
					str, err := dom.GetOuterHTML().WithNodeID(v.NodeID).Do(c)
					if err != nil {
						return err
					}
					res.WriteString(str)
				}
				return nil
			}),
		)
	} else {
		tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			str, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				return err
			}
			res.WriteString(str)
			return nil
		}))
	}
	if c.screenshot {
		tasks = append(tasks, chromedp.CaptureScreenshot(buf))
	}
	return tasks
}
