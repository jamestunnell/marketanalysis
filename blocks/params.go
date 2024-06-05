package blocks

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/exp/maps"
)

type Params map[string]Param
type ParamVals map[string]any

type errInvalidParam struct {
	Name   string
	Errors []error
}

func (ps Params) GetSortedNames() []string {
	names := maps.Keys(ps)

	sort.Strings(names)

	return names
}

func (ps Params) GetNonDefaultValues() ParamVals {
	vals := ParamVals{}

	for name, ps := range ps {
		if val := ps.GetCurrentVal(); val != ps.GetDefaultVal() {
			vals[name] = val
		}
	}

	return vals
}

func (ps Params) SetValuesOrDefault(vals ParamVals) error {
	errs := []error{}

	for name, p := range ps {
		val, found := vals[name]
		if !found {
			val = p.GetDefaultVal()
		}

		if err := p.SetCurrentVal(val); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		var merr *multierror.Error

		for _, err := range errs {
			merr = multierror.Append(merr, err)
		}

		return fmt.Errorf("failed to set param values: %w", merr)
	}

	return nil
}
