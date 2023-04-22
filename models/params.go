package models

import (
	"encoding/json"
	"sort"

	"golang.org/x/exp/maps"
)

type Params map[string]Param

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
