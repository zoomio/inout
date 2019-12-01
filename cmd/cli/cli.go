package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/zoomio/inout"
)

func main() {
	source := flag.String("s", "", "Source, e.g. \"-s https://stackoverflow.com\"")
	query := flag.String("q", "", "DOM CSS query, waits until element available, "+
		"e.g. \"-q p\" will fetch contents of all <p> tags on the given source")
	timeout := flag.Duration("t", 5*time.Second, "\"-d 1s\" sets timeout to 1 second")
	verbose := flag.Bool("v", false, "\"-v\" enables verbose mode")
	flag.Parse()

	reader, err := inout.NewInOut(
		context.Background(),
		inout.Source(*source),
		inout.Query(*query),
		inout.Timeout(*timeout),
		inout.Verbose(*verbose))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	lines, err := reader.ReadLines()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	fmt.Println(lines)
}
