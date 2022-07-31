package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/zoomio/inout"
)

var (
	source  = flag.String("s", "", "Source, e.g. \"-s https://stackoverflow.com\"")
	query   = flag.String("q", "", "DOM CSS query, waits until element available, e.g. \"-q p\" will return contents of all <p> tags on the given source")
	ready   = flag.String("r", "", "DOM CSS query, waits until element available, returns the whole HTML document")
	until   = flag.Duration("u", 0, "duration to wait before getting HTML contents, handy for SPAs, because they keep loading in browsers for some time")
	img     = flag.String("i", "", "enables capturing screenshot in the provided path")
	timeout = flag.Duration("t", 5*time.Second, "timeout for the whole fetch, e.g. \"-d 1s\" sets timeout to 1 second")
	verbose = flag.Bool("v", false, "\"-v\" enables verbose mode")
)

func main() {
	flag.Parse()

	ttl := *timeout
	if ttl > 0 && *until > 0 {
		ttl += *until
	}

	reader, err := inout.NewInOut(
		context.Background(),
		inout.Source(*source),
		inout.Query(*query),
		inout.WaitFor(*ready),
		inout.WaitUntil(*until),
		inout.Screenshot(len(*img) > 0),
		inout.Timeout(ttl),
		inout.Verbose(*verbose))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create reader: %v\n", err)
		os.Exit(1)
	}
	defer reader.Close()

	if len(*img) > 0 && len(reader.ImgBytes) > 0 {
		err = ioutil.WriteFile(*img, reader.ImgBytes, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to store captured image at %s: %v\n", *img, err)
			os.Exit(3)
		}
	}

	for {
		line, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "failed to read line: %v\n", err)
			os.Exit(2)
		}
		fmt.Println(line)
	}
}
