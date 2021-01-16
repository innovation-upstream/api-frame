package repo

import (
	"cloud.google.com/go/firestore"
)

type firestoreQueryMutate struct {
	query *firestore.Query
	limit int
}

func (m *firestoreQueryMutate) ApplyMutations(opts ...QueryOption) {
	for _, opt := range opts {
		opt(m)
	}
}

func (m *firestoreQueryMutate) SetLimit(l int) {
	m.limit = l
	// Add one to limit so we know if there are more results
	m.query.Limit(l + 1)
}

func (m *firestoreQueryMutate) SetOffset(o int) {
	m.query.Offset(o)
}

func (m *firestoreQueryMutate) GetLimit() int {
	return m.limit
}
