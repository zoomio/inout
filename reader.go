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
	ImgBytes []byte

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
	cfg := &config{}

	// apply custom configuration
	for _, option := range options {
		option(cfg)
	}

	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	start := time.Now()
	end := start.Add(cfg.timeout)
	go func() {
		ticker := time.NewTicker(30 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if cfg.timeout > 0 && end.Before(time.Now()) {
					if cfg.verbose {
						fmt.Printf("timeout of %v passed, stopping...\n", cfg.timeout)
					}
					cancel()
					return
				}
			}
		}
	}()

	r := &Reader{}
	var err error

	// STDIN
	if cfg.source == "" {
		r.reader, err = handleSTDIN()
		if err != nil {
			return *r, err
		}

		// HTTP
	} else if cfg.isHTTP() {
		if cfg.verbose {
			fmt.Println("source is HTTP/HTTPS")
		}
		r.reader, r.ImgBytes, err = handleHTTP(childCtx, cfg)
		if err != nil {
			return *r, err
		}

		// File system
	} else if cfg.isFS() {
		r.reader, err = handleFS(cfg.source)
		if err != nil {
			return *r, err
		}

		// Unresolvable "source"
	} else {
		return *r, fmt.Errorf("unknown type of provided input source: %s", cfg.source)
	}

	r.buffer = make([]byte, 64*1024)
	r.scanner = bufio.NewScanner(r.reader)
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

type config struct {
	source     string
	query      string
	waitFor    string
	waitUntil  time.Duration
	screenshot bool
	timeout    time.Duration
	verbose    bool
	userAgent  string
}

func (c *config) isHTTP() bool {
	return strings.HasPrefix(c.source, "http") || strings.HasPrefix(c.source, "https")
}

func (c *config) isHeadless() bool {
	return len(c.query) > 0 || len(c.waitFor) > 0 || c.waitUntil > 0 || c.screenshot
}

func (c *config) isFS() bool {
	_, err := os.Stat(c.source)
	return err == nil
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

func handleHTTP(ctx context.Context, c *config) (io.ReadCloser, []byte, error) {
	if c.isHeadless() {
		res, err := headless(ctx, c)
		if err != nil {
			return nil, nil, fmt.Errorf("error in headless query %#v: %w", c, err)
		}
		return io.NopCloser(strings.NewReader(res.htmlDoc)), res.imgBytes, nil
	}
	res := fetch(c.source)
	if res.err != nil {
		return nil, nil, fmt.Errorf("error in fetching provided source=%s: %w", c.source, res.err)
	}
	return res.resp.Body, nil, nil
}

func handleFS(source string) (io.Reader, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("error in opening file source=%s for reading: %w", source, err)
	}
	return bufio.NewReader(f), nil
}
