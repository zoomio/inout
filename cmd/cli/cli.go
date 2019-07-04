package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zoomio/inout"
)

func main() {
	source := flag.String("s", "", "Source, e.g. -s https://stackoverflow.com")
	query := flag.String("q", "", "DOM CSS query, waits until element available, "+
		"e.g. `-q p` will fetch contents of all <p> tags on the given source")
	flag.Parse()

	reader, err := inout.New(*source, *query)
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
