package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zoomio/inout"
)

func main() {
	source := flag.String("s", "", "Source, e.g. -s https://stackoverflow.com")
	selector := flag.String("dom", "", "DOM CSS selector, waits until element available, e.g. -dom 'p'")
	flag.Parse()

	reader, err := inout.New(*source, *selector)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	lines, err := reader.ReadLines()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}

	fmt.Println(lines)
}
