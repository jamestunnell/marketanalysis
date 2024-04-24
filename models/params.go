package models

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/exp/maps"
)

type Params map[string]Param

type errInvalidParam struct {
	Name   string
	Errors []error
}

func (ps Params) SortedNames() []string {
	names := maps.Keys(ps)

	sort.Strings(names)

	return names
}

func (ps Params) String() string {
	paramValsByName := map[string]any{}

	for name, ps := range ps {
		paramValsByName[name] = ps.GetVal()
	}

	d, _ := json.Marshal(paramValsByName)

	return string(d)
}

func (ps Params) Validate() error {
	errs := []error{}

	for name, p := range ps {
		val := p.GetVal()

		paramErrs := []error{}

		for _, c := range p.Constraints() {
			if paramErr := c.Check(val); paramErr != nil {
				paramErrs = append(paramErrs, paramErr)
			}
		}

		if len(paramErrs) > 0 {
			errs = append(errs, &errInvalidParam{Name: name, Errors: paramErrs})
		}
	}

	if len(errs) > 0 {
		var merr *multierror.Error

		for _, err := range errs {
			merr = multierror.Append(merr, err)
		}

		return merr
	}

	return nil
}

func (err *errInvalidParam) Error() string {
	var merr *multierror.Error

	for _, err := range err.Errors {
		merr = multierror.Append(merr, err)
	}

	return fmt.Sprintf("invalid param %s: %v", err.Name, merr)
}
