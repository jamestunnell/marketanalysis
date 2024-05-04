package optimization

import (
	"math/rand"
	"reflect"

	"github.com/jamestunnell/marketanalysis/blocks"
)

func MakePartial(param blocks.Param, rng *rand.Rand) (PartialGenome, error) {
	switch param.GetType() {
	case "int":
		return makeIntPartial(param, rng)
	case "float64":
		return makeFloatPartial(param, rng)
	}

	return nil, &ErrUnsupportedType{Type: param.GetType()}
}

func makeIntPartial(param blocks.Param, rng *rand.Rand) (PartialGenome, error) {
	var m IntMutator

	switch p := param.(type) {
	case *blocks.IntRange:
		m = NewIntRangeMutator(p.Min, p.Max)
	case *blocks.IntEnum:
		m = NewIntEnumMutator(p.Enum)
	}

	if m == nil {
		return nil, NewErrUnsupportedCombo("int", reflect.TypeOf(param).String())
	}

	return NewIntValue(m, rng), nil
}

func makeFloatPartial(param blocks.Param, rng *rand.Rand) (PartialGenome, error) {
	var m FloatMutator

	switch p := param.(type) {
	case *blocks.FltRange:
		m = NewFloatRangeMutator(p.Min, p.Max)
	case *blocks.FltEnum:
		m = NewFloatEnumMutator(p.Enum)
	}

	if m == nil {
		return nil, NewErrUnsupportedCombo("float64", reflect.TypeOf(param).String())
	}

	return NewFloatValue(m, rng), nil
}
