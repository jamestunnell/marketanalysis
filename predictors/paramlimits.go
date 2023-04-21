package predictors

import "golang.org/x/exp/maps"

var upperLimits = map[string]UpperLimit{}

func UpperLimitNames() []string {
	return maps.Keys(upperLimits)
}

func GetUpperLimit(name string) (UpperLimit, bool) {
	limit, found := upperLimits[name]
	if !found {
		return nil, false
	}

	return limit, true
}
