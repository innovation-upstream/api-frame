package repo

type QueryMutate interface {
	ApplyMutations(opts ...QueryOption)
	GetLimit() int
	SetLimit(l int)
	SetOffset(o int)
}
