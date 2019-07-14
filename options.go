package inout

// Option allows to customise instance of In-Out.
type Option func(*Reader)

var (
	// WithSource sets target source of the In-Out.
	WithSource = func(source string) Option {
		return func(r *Reader) {
			r.source = source
		}
	}

	// WithQuery sets CSS query for the target of In-Out.
	WithQuery = func(query string) Option {
		return func(r *Reader) {
			r.query = query
		}
	}
)
