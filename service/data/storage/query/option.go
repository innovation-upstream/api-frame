package query

// Option is for customizing queries
type Option func(Customize)

// WithLimit sets the query limit
func WithLimit(l int) Option {
	return func(q Customize) {
		q.SetLimit(l)
	}
}

// WithLimit sets the query limit
func WithOffset(l int) Option {
	return func(q Customize) {
		q.SetOffset(l)
	}
}

// WithFields sets the fields that are to be set or retrieved
func WithFields(f []string) Option {
	return func(q Customize) {
		q.SetQueryFields(f)
	}
}

// WithWhere adds a where clause
func WithWhere(path, op string, val interface{}) Option {
	return func(q Customize) {
		q.AddWhere(path, op, val)
	}
}
