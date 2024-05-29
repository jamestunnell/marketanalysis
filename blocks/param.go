package blocks

type Param interface {
	GetType() string
	GetDefault() any
	GetLimits() []any
	GetVal() any
	SetVal(any) error
}
