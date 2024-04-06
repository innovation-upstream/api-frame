package query

import (
	"cloud.google.com/go/firestore"
	"github.com/innovation-upstream/api-frame/service/data/storage/query"
)

type (
	firestoreQueryCustomize struct {
		query  *firestore.Query
		limit  int
		fields []string
	}

	FirestoreQueryCustomizeFactory func(q *firestore.Query) query.Customize
)

// NewFirestoreQueryCustomize constructs a customizeFirestore instance
var NewFirestoreQueryCustomize = FirestoreQueryCustomizeFactory(func(q *firestore.Query) query.Customize {
	return &firestoreQueryCustomize{
		query: q,
	}
})

func (c *firestoreQueryCustomize) ApplyOptions(opts ...query.Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *firestoreQueryCustomize) SetLimit(l int) {
	c.limit = l
	// Add one to limit so we know if there are more results
	queryWithLimit := c.query.Limit(l + 1)
	*c.query = queryWithLimit
}

func (c *firestoreQueryCustomize) SetOffset(o int) {
	queryWithOffset := c.query.Offset(o)
	*c.query = queryWithOffset
}

func (c *firestoreQueryCustomize) GetLimit() int {
	return c.limit
}

func (c *firestoreQueryCustomize) SetQueryFields(f []string) {
	c.fields = f
}

func (c *firestoreQueryCustomize) GetQueryFields() []string {
	return c.fields
}

func (c *firestoreQueryCustomize) AddWhere(path, op string, val interface{}) {
	queryWithWhere := c.query.Where(path, op, val)
	*c.query = queryWithWhere
}

func (c *firestoreQueryCustomize) SetStartAfter(field string, val interface{}) {
	queryWithStart := c.query.OrderBy(field, firestore.Asc).StartAfter(val)
	*c.query = queryWithStart
}
