package params

type Int struct {
	value int
}

const TypeInt = "Int"

func NewInt(value int) *Int {
	return &Int{value: value}
}

func (p *Int) Type() string {
	return TypeInt
}

func (p *Int) Value() any {
	return p.value
}
