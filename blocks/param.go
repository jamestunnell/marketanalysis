package blocks

type Param interface {
	GetType() string
	GetDefault() any
	GetSchema() map[string]any
	GetVal() any
	SetVal(any) error
}
