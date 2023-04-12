package params

type Float struct {
	value float64
}

const TypeFloat = "Float"

func NewFloat(value float64) *Float {
	return &Float{value: value}
}

func (p *Float) Type() string {
	return TypeFloat
}

func (p *Float) Value() any {
	return p.value
}
