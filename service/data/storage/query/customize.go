package query

// Customize is for decoupling external query interfaces from storage consumers
type Customize interface {
	ApplyOptions(opts ...Option)
	GetLimit() int
	SetLimit(l int)
	SetOffset(o int)
	SetQueryFields(fields []string)
	GetQueryFields() []string
	AddWhere(field, op string, val interface{})
}
