package inout

import "time"

// Option allows to customise instance of In-Out.
type Option func(*Reader)

var (
	// Source sets target source of the In-Out.
	Source = func(source string) Option {
		return func(r *Reader) {
			r.source = source
		}
	}

	// Query sets CSS query for the target of In-Out.
	Query = func(query string) Option {
		return func(r *Reader) {
			r.query = query
		}
	}

	// Timeout sets timeout for the operation.
	Timeout = func(timeout time.Duration) Option {
		return func(r *Reader) {
			r.timeout = timeout
		}
	}

	// Verbose enables verbose mode.
	Verbose = func(verbose bool) Option {
		return func(r *Reader) {
			r.verbose = verbose
		}
	}
)
