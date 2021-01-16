package repo

type QueryOption func(QueryMutate)

// WithLimit sets the query limit
func WithLimit(l int) QueryOption {
	return func(q QueryMutate) {
		q.SetLimit(l)
	}
}

// WithLimit sets the query limit
func WithOffset(l int) QueryOption {
	return func(q QueryMutate) {
		q.SetOffset(l)
	}
}
