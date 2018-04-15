package inout

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// In - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type In struct {
	reader io.Reader
}

// NewIn initializes an input from STDIN, file or web page.
//
// name - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func NewIn(name string) In {
	var reader io.Reader

	// STDIN
	if name == "" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			panic(fmt.Sprintf("error in reading from STDIN: %v", err))
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return In{}
		}
		reader = bufio.NewReader(os.Stdin)
		// File system
	} else if _, err := os.Stat(name); err == nil {
		f, err := os.Open(name)
		if err != nil {
			panic(fmt.Sprintf("error in opening file %s for reading: %v", name, err))
		}
		reader = bufio.NewReader(f)
		// HTTP
	} else {
		resp, err := http.Get(name)
		if err != nil {
			panic(fmt.Sprintf("provided name=%s is not a file and not a URL: %v", name, err))
		}
		reader = resp.Body
	}

	return In{
		reader: reader,
	}
}

// NewInFromString initializes an input from string.
func NewInFromString(input string) In {
	return In{
		reader: strings.NewReader(input),
	}
}

// ReadString ...
func (in *In) ReadString() string {
	s := bufio.NewScanner(in.reader)
	var text string
	if s.Scan() {
		text = s.Text()
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
	return text
}

// ReadInt ...
func (in *In) ReadInt() int {
	i, err := strconv.ParseInt(in.ReadString(), 10, 64)
	if err != nil {
		panic(err)
	}
	return int(i)
}

// ReadAllStrings provides slice of strings from input split by white space.
func (in *In) ReadAllStrings() []string {
	tokens := make([]string, 0)
	lines, err := in.lines()
	if err != nil {
		panic(fmt.Sprintf("error in reading from scanner: %v", err))
	}
	for _, line := range lines {
		tokens = append(tokens, strings.Fields(line)...)
	}
	return tokens
}

// ReadAllInts reads all remaining tokens from this input stream, parses them as integers,
// and returns them as an array of integers.
//
// Returns all remaining lines in this input stream, as an array of integers
func (in *In) ReadAllInts() []int {
	fields := in.ReadAllStrings()
	vals := make([]int, len(fields))
	for i, f := range fields {
		n, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("error in parsing %s: %v", f, err))
		}
		vals[i] = int(n)
	}
	return vals
}

// Close ...
func (in *In) Close() {
	if closer, ok := in.reader.(io.Closer); ok {
		defer closer.Close()
	}
}

func (in *In) lines() ([]string, error) {
	defer in.Close()
	s := bufio.NewScanner(in.reader)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
