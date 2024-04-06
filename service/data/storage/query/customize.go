package query

//go:generate mockgen -destination=../mock/customize_mock.go -package=mock github.com/innovation-upstream/api-frame/service/data/storage/query Customize

// Customize is for decoupling external query interfaces from storage consumers
type Customize interface {
	ApplyOptions(opts ...Option)
	GetLimit() int
	SetLimit(l int)
	SetOffset(o int)
	SetQueryFields(fields []string)
	GetQueryFields() []string
	AddWhere(field, op string, val interface{})
	SetStartAfter(field string, val interface{})
}
