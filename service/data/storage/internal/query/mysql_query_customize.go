package query

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.innovationup.stream/innovation-upstream/api-frame/service/data/storage/query"
)

type mysqlQueryCustomize struct {
	selectBuilder *sq.SelectBuilder
	updateBuilder *sq.UpdateBuilder
	insertBuilder *sq.InsertBuilder
	limit         int
	fields        []string
}

// NewMysqlQueryCustomize constructs a mysqlQueryCustomize instance
func NewMysqlQueryCustomize(selectBuilder *sq.SelectBuilder, updateBuilder *sq.UpdateBuilder, insertBuilder *sq.InsertBuilder) query.Customize {
	return &mysqlQueryCustomize{
		selectBuilder: selectBuilder,
		updateBuilder: updateBuilder,
		insertBuilder: insertBuilder,
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
	} else {
		c.updateBuilder.Limit(uint64(l + 1))
	}
}

func (c *mysqlQueryCustomize) SetOffset(o int) {
	if c.selectBuilder != nil {
		c.selectBuilder.Offset(uint64(o + 1))
	} else {
		c.updateBuilder.Offset(uint64(o + 1))
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
	} else {
		c.updateBuilder.Where(path, op, val)
	}
}
