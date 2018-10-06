package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zoomio/inout"
)

func main() {
	source := flag.String("s", "", "Source")
	flag.Parse()

	reader, err := inout.New(*source)
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
