package appvars

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Vars struct {
	Debug  bool
	Port   int
	DBConn string
	// DBUser, DBPass string
}

type VarCandidate struct {
	Source string
	Value  string
}

const (
	DefaultPort  = "4002"
	DefaultDebug = "false"
)

var errVarNotFound = errors.New("var not found")

func loadAppVar[T comparable](
	name string,
	parse func(string) (T, error),
	first *VarCandidate,
	more ...*VarCandidate) (T, error) {
	allSources := []string{}
	candidates := append([]*VarCandidate{first}, more...)

	var source string
	var valStr string

	for _, c := range candidates {
		if c.Value != "" {
			valStr = c.Value
			source = c.Source

			break
		}

		allSources = append(allSources, c.Source)
	}

	if valStr == "" {
		var val T

		return val, errVarNotFound
	}

	value, err := parse(valStr)
	if err != nil {
		var val T

		return val, fmt.Errorf("failed to parse '%s': %w", valStr, err)
	}

	log.Info().
		Str("name", name).
		Str("source", source).
		Interface("value", value).
		Msgf("loaded app var")

	return value, nil
}

func newVarCandidate(val, source string) *VarCandidate {
	return &VarCandidate{
		Source: source,
		Value:  val,
	}
}
