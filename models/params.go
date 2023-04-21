package models

import "encoding/json"

type Params map[string]Param

func (ps Params) String() string {
	paramValsByName := map[string]any{}

	for name, ps := range ps {
		paramValsByName[name] = ps.GetVal()
	}

	d, _ := json.Marshal(paramValsByName)

	return string(d)
}
