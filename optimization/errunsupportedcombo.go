package optimization

import "fmt"

type ErrUnsupportedCombo struct {
	ValueType, ConstraintType string
}

func NewErrUnsupportedCombo(valType, constraintType string) *ErrUnsupportedCombo {
	return &ErrUnsupportedCombo{
		ValueType:      valType,
		ConstraintType: constraintType,
	}
}
func (err *ErrUnsupportedCombo) Error() string {
	const strFmt = "unsupported combo: value type '%s' with constraint type '%s'"

	return fmt.Sprintf(strFmt, err.ValueType, err.ConstraintType)
}
