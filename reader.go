package inout

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Reader - Input. This struct provides methods for reading strings
// and numbers from standard input, file input, URLs, and sockets.
type Reader struct {
	source  string
	query   string
	timeout time.Duration
	verbose bool
	reader  io.Reader
	buffer  []byte
	scanner *bufio.Scanner
}

// New initializes an instance of Reader from STDIN, file or web page.
//
// source - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func New(ctx context.Context, source string) (Reader, error) {
	return NewInOut(ctx, Source(source))
}

// NewInOut initializes an instance of Reader from STDIN, file or web page.
//
// source - the filename or web page name, reads from STDIN if name is empty.
// Panics on errors.
func NewInOut(ctx context.Context, options ...Option) (Reader, error) {
	var reader io.Reader
	var err error

	r := &Reader{}

	// apply custom configuration
	for _, option := range options {
		option(r)
	}

	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	start := time.Now()
	end := start.Add(r.timeout)
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if r.timeout > 0 && end.Before(time.Now()) {
					if r.verbose {
						fmt.Printf("timeout of %v passed, stopping...\n", r.timeout)
					}
					cancel()
					return
				}
			}
		}
	}()

	// STDIN
	if r.source == "" {
		reader, err = handleSTDIN()
		if err != nil {
			return *r, err
		}

		// HTTP
	} else if strings.HasPrefix(r.source, "http") || strings.HasPrefix(r.source, "https") {
		if r.verbose {
			fmt.Println("source is HTTP/HTTPS")
		}
		reader, err = handleHTTP(childCtx, r.source, r.query, r.verbose)
		if err != nil {
			return *r, err
		}

		// File system
	} else if _, err := os.Stat(r.source); err == nil {
		reader, err = handleFS(r.source)
		if err != nil {
			return *r, err
		}

		// Unresolvable "source"
	} else {
		return *r, fmt.Errorf("unknown type of provided input source: %s", r.source)
	}

	r.reader = reader
	r.buffer = make([]byte, 64*1024)
	r.scanner = bufio.NewScanner(reader)
	r.scanner.Buffer(r.buffer, 1024*1024)

	return *r, nil
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
	} else {
		return "", io.EOF
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
		return nil, fmt.Errorf("error in reading from STDIN: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, errors.New("unsupported mode")
	}
	return bufio.NewReader(os.Stdin), nil
}

func handleHTTP(ctx context.Context, source, query string, verbose bool) (io.Reader, error) {
	if query != "" {
		text, err := waitForDomElement(ctx, query, source, verbose)
		if err != nil {
			return nil, fmt.Errorf("error in waiting for query=%s in source=%s: %w", query, source, err)
		}
		return strings.NewReader(text), nil
	}
	res := fetch(ctx, source)
	if res.err != nil {
		return nil, fmt.Errorf("error in fetching provided source=%s: %w", source, res.err)
	}
	return res.resp.Body, nil
}

func handleFS(source string) (io.Reader, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("error in opening file source=%s for reading: %w", source, err)
	}
	return bufio.NewReader(f), nil
}
