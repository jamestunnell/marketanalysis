package optimization

import (
	"math/rand"
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/models"
)

func MakePartial(param models.Param, rng *rand.Rand) (PartialGenome, error) {
	switch param.Type() {
	case "int":
		return makeIntPartial(param, rng)
	case "float64":
		return makeFloatPartial(param, rng)
	}

	return nil, &ErrUnsupportedType{Type: param.Type()}
}

func makeIntPartial(param models.Param, rng *rand.Rand) (PartialGenome, error) {
	cs := param.Constraints()
	if len(cs) != 1 {
		return nil, commonerrs.NewErrExactLen("constraints", len(cs), 1)
	}

	var m IntMutator

	switch c := cs[0].(type) {
	case *constraints.TypedValRange[int]:
		m = NewIntRangeMutator(c.Min, c.Max)
	case *constraints.TypedValOneOf[int]:
		m = NewIntEnumMutator(c.Allowed)
	}

	if m == nil {
		return nil, NewErrUnsupportedCombo("int", reflect.TypeOf(cs[0]).String())
	}

	return NewIntValue(m, rng), nil
}

func makeFloatPartial(param models.Param, rng *rand.Rand) (PartialGenome, error) {
	cs := param.Constraints()
	if len(cs) != 1 {
		return nil, commonerrs.NewErrExactLen("constraints", len(cs), 1)
	}

	var m FloatMutator

	switch c := cs[0].(type) {
	case *constraints.TypedValRange[float64]:
		m = NewFloatRangeMutator(c.Min, c.Max)
	case *constraints.TypedValOneOf[float64]:
		m = NewFloatEnumMutator(c.Allowed)
	}

	if m == nil {
		return nil, NewErrUnsupportedCombo("float64", reflect.TypeOf(cs[0]).String())
	}

	return NewFloatValue(m, rng), nil
}
