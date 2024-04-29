package blocks

type Param interface {
	GetDefault() any
	GetSchema() map[string]any
	GetVal() any
	SetVal(any) error
}
