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
func New(source, query string) (*Reader, error) {
	var reader io.Reader
	var err error

	// STDIN
	if source == "" {
		reader, err = handleSTDIN()
		if err != nil {
			return nil, err
		}

		// HTTP
	} else if strings.HasPrefix(source, "http") || strings.HasPrefix(source, "https") {
		reader, err = handleHTTP(source, query)
		if err != nil {
			return nil, err
		}

		// File system
	} else if _, err := os.Stat(source); err == nil {
		reader, err = handleFS(source)
		if err != nil {
			return nil, err
		}

		// Unresolvable "source"
	} else {
		return nil, fmt.Errorf("unknown type of provided input source: %s", source)
	}

	return &Reader{
		reader:  reader,
		scanner: bufio.NewScanner(reader),
	}, nil
}

// NewFromString initializes an input from string.
func NewFromString(input string) *Reader {
	r := strings.NewReader(input)
	return &Reader{
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

func handleSTDIN() (io.Reader, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("error in reading from STDIN: %v", err)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, errors.New("unsupported mode")
	}
	return bufio.NewReader(os.Stdin), nil
}

func handleHTTP(source, query string) (io.Reader, error) {
	if query != "" {
		text, err := waitForDomElement(query, source)
		if err != nil {
			return nil, fmt.Errorf("error in waiting for %s in %s: %v", query, source, err)
		}
		return strings.NewReader(text), nil
	}
	resp, err := http.Get(source)
	if err != nil {
		return nil, fmt.Errorf("provided source=%s is not a file and not a URL: %v", source, err)
	}
	return resp.Body, nil
}

func handleFS(source string) (io.Reader, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("error in opening file %s for reading: %v", source, err)
	}
	return bufio.NewReader(f), nil
}
