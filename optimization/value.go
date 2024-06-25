package optimization

import (
	"math/rand"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type Value interface {
	Init(rng *rand.Rand)
	GetValue() any
	Mutate(rng *rand.Rand)
	Clone() Value
}

type Values map[string]Value

func MakeValue(constraint *models.ConstraintInfo) (Value, error) {
	switch constraint.Type {
	case "OneOf[int]", "RangeIncl[int]":
		intLimits := sliceutils.Map(constraint.Limits, func(val any) int { return int(val.(float64)) })

		var mut IntMutator

		if constraint.Type == "OneOf[int]" {
			mut = NewIntEnumMutator(intLimits)
		} else if intLimits[0] == intLimits[1] {
			mut = NewIntConstMutator(intLimits[0])
		} else {
			mut = NewIntRangeMutator(intLimits[0], intLimits[1])
		}
		return NewIntValue(mut), nil
	case "OneOf[float64]", "RangeIncl[float64]":
		fltLimits := sliceutils.Map(constraint.Limits, func(val any) float64 { return val.(float64) })

		var mut FloatMutator

		if constraint.Type == "OneOf[float64]" {
			mut = NewFloatEnumMutator(fltLimits)
		} else if fltLimits[0] == fltLimits[1] {
			mut = NewFloatConstMutator(fltLimits[0])
		} else {
			mut = NewFloatRangeMutator(fltLimits[0], fltLimits[1])
		}

		return NewFloatValue(mut), nil
	}

	return nil, &ErrUnsupportedType{Type: constraint.Type}
}
