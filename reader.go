package inout

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Reader - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type Reader struct {
	reader  io.Reader
	scanner *bufio.Scanner
}

// New initializes an instance of Reader from STDIN, file or web page.
//
// source - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func New(source, selector string) (Reader, error) {
	var reader io.Reader
	in := Reader{}

	// STDIN
	if source == "" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return in, fmt.Errorf("error in reading from STDIN: %v", err)
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return in, errors.New("unsupported mode")
		}
		reader = bufio.NewReader(os.Stdin)

		// HTTP
	} else if strings.HasPrefix(source, "http") || strings.HasPrefix(source, "https") {

		if (selector != "") {
			text, err := waitForDomElement(selector, source)
			if err != nil {
				return in, fmt.Errorf("error in waiting for %s in %s: %v", selector, source, err)
			}
			reader = strings.NewReader(text)
		} else {
			resp, err := http.Get(source)
			if err != nil {
				return in, fmt.Errorf("provided source=%s is not a file and not a URL: %v", source, err)
			}
			reader = resp.Body
		}

		// File system
	} else if _, err := os.Stat(source); err == nil {
		f, err := os.Open(source)
		if err != nil {
			return in, fmt.Errorf("error in opening file %s for reading: %v", source, err)
		}
		reader = bufio.NewReader(f)

		// Unresolvable "source"
	} else {
		return Reader{}, fmt.Errorf("unknown type of provided input source: %s", source)
	}

	return Reader{
		reader:  reader,
		scanner: bufio.NewScanner(reader),
	}, nil
}

// NewFromString initializes an input from string.
func NewFromString(input string) Reader {
	r := strings.NewReader(input)
	return Reader{
		reader:  r,
		scanner: bufio.NewScanner(r),
	}
}

// Read reads into given bytes (does not close reader).
func (in *Reader) Read(p []byte) (n int, err error) {
	return in.reader.Read(p)
}

// ReadLine reads line from reader (does not close reader).
func (in *Reader) ReadLine() (string, error) {
	var text string
	if in.scanner.Scan() {
		text = in.scanner.Text()
	}
	err := in.scanner.Err()
	if err != nil {
		return "", err
	}
	return text, nil
}

// ReadWords provides slice of all words from input split by white space and closes the reader.
func (in *Reader) ReadWords() ([]string, error) {
	tokens := make([]string, 0)
	lines, err := in.ReadLines()
	if err != nil {
		return tokens, fmt.Errorf("error in reading from scanner: %v", err)
	}
	for _, line := range lines {
		tokens = append(tokens, strings.Fields(line)...)
	}
	return tokens, nil
}

// Close closes reader.
func (in *Reader) Close() error {
	if closer, ok := in.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// ReadLines provides slice of all text lines from input and closes the reader.
func (in *Reader) ReadLines() ([]string, error) {
	defer in.Close()
	var lines []string
	for in.scanner.Scan() {
		lines = append(lines, in.scanner.Text())
	}
	if err := in.scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
