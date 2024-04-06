package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/innovation-upstream/api-frame/service/data/storage/query"
)

type mysqlQueryCustomize struct {
	selectBuilder *sq.SelectBuilder
	updateBuilder *sq.UpdateBuilder
	insertBuilder *sq.InsertBuilder
	deleteBuilder *sq.DeleteBuilder
	limit         int
	fields        []string
}

// NewMysqlQueryCustomize constructs a mysqlQueryCustomize instance
func NewMysqlQueryCustomize(selectBuilder *sq.SelectBuilder, updateBuilder *sq.UpdateBuilder, insertBuilder *sq.InsertBuilder, deleteBuilder *sq.DeleteBuilder) query.Customize {
	return &mysqlQueryCustomize{
		selectBuilder: selectBuilder,
		updateBuilder: updateBuilder,
		insertBuilder: insertBuilder,
		deleteBuilder: deleteBuilder,
	}
}

func (c *mysqlQueryCustomize) ApplyOptions(opts ...query.Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *mysqlQueryCustomize) SetLimit(l int) {
	c.limit = l
	if c.selectBuilder != nil {
		c.selectBuilder.Limit(uint64(l + 1))
	} else if c.updateBuilder != nil {
		c.updateBuilder.Limit(uint64(l + 1))
	} else if c.deleteBuilder != nil {
		c.deleteBuilder.Limit(uint64(l + 1))
	}
}

func (c *mysqlQueryCustomize) SetOffset(o int) {
	if c.selectBuilder != nil {
		c.selectBuilder.Offset(uint64(o + 1))
	} else if c.updateBuilder != nil {
		c.updateBuilder.Offset(uint64(o + 1))
	} else if c.deleteBuilder != nil {
		c.deleteBuilder.Offset(uint64(o + 1))
	}
}

func (c *mysqlQueryCustomize) GetLimit() int {
	return c.limit
}

func (c *mysqlQueryCustomize) SetQueryFields(f []string) {
	c.fields = f
}

func (c *mysqlQueryCustomize) GetQueryFields() []string {
	return c.fields
}

func (c *mysqlQueryCustomize) AddWhere(path, op string, val interface{}) {
	if c.selectBuilder != nil {
		c.selectBuilder.Where(path, op, val)
	} else if c.updateBuilder != nil {
		c.updateBuilder.Where(path, op, val)
	} else if c.deleteBuilder != nil {
		c.deleteBuilder.Where(path, op, val)
	}
}

// SetStartAfter is a no-op for mysql
func (c *mysqlQueryCustomize) SetStartAfter(field string, val interface{}) {
	fmt.Println("SetStartAfter is a no-op for mysql")
}
