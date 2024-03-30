package inout

import "time"

// Option allows to customise instance of In-Out.
type Option func(*config)

var (
	// Source sets target source of the In-Out.
	Source = func(source string) Option {
		return func(c *config) {
			c.source = source
		}
	}

	// Query sets CSS query for the target of In-Out.
	Query = func(query string) Option {
		return func(c *config) {
			c.query = query
		}
	}

	// WaitFor sets CSS query for the target of In-Out.
	WaitFor = func(query string) Option {
		return func(c *config) {
			c.waitFor = query
		}
	}

	// WaitUntil sets page load duration to wait for.
	WaitUntil = func(d time.Duration) Option {
		return func(c *config) {
			c.waitUntil = d
		}
	}

	// Screenshot captures screenshot, Reader will ImgBytes of the image populated.
	Screenshot = func(v bool) Option {
		return func(c *config) {
			c.screenshot = v
		}
	}

	// Timeout sets timeout for the operation.
	Timeout = func(timeout time.Duration) Option {
		return func(c *config) {
			c.timeout = timeout
		}
	}

	// Verbose enables verbose mode.
	Verbose = func(verbose bool) Option {
		return func(c *config) {
			c.verbose = verbose
		}
	}

	// UserAgent custom user agent for healess Chrome operations.
	UserAgent = func(ua string) Option {
		return func(c *config) {
			c.userAgent = ua
		}
	}
)
