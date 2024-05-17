package graph

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/blocks/registry"
)

const ConfigKeyName = "id"

type Configuration struct {
	ID          string         `json:"id" bson:"_id"`
	Name        string         `json:"name"`
	Blocks      []*BlockConfig `json:"blocks"`
	Connections []*Connection  `json:"connections"`
}

type BlockConfig struct {
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	ParamVals map[string]any `json:"paramVals"`
	Recording []string       `json:"recording"`
}

func (cfg Configuration) GetKey() string {
	return cfg.ID
}

func (cfg *Configuration) Validate() []error {
	errs := []error{}
	blks := Blocks{}

	for _, b := range cfg.Blocks {
		if _, found := blks[b.Name]; found {
			err := fmt.Errorf("duplicate block name '%s'", b.Name)

			errs = append(errs, err)

			continue
		}

		newBlk, blkFound := registry.Get(b.Type)
		if !blkFound {
			err := fmt.Errorf("block %s: has unknown type '%s'", b.Name, b.Type)

			errs = append(errs, err)

			continue
		}

		blk := newBlk()

		blks[b.Name] = blk

		outs := blk.GetOutputs()

		// validate recording outputs
		for _, recOut := range b.Recording {
			if out, found := outs[recOut]; !found {
				errs = append(errs, fmt.Errorf("block %s: recording output '%s' not found", b.Name, recOut))
			} else if out.GetType() != "float64" {
				errs = append(errs, fmt.Errorf("block %s: recording output '%s' not a float64 type", b.Name, recOut))
			}
		}

		// validate params
		params := blk.GetParams()

		for name, val := range b.ParamVals {
			param, found := params[name]
			if !found {
				err := fmt.Errorf("block %s: unknown param name '%s'", b.Name, name)

				errs = append(errs, err)

				continue
			}

			if err := param.SetVal(val); err != nil {
				err = fmt.Errorf("block %s: param %s: value %v is invalid: %w", b.Name, name, val, err)

				errs = append(errs, err)
			}
		}
	}

	// validate connections
	for i, conn := range cfg.Connections {
		if conn == nil {
			err := fmt.Errorf("connection #%d is null", i+1)

			errs = append(errs, err)

			continue
		}

		if _, found := blks.FindOutput(conn.Source); !found {
			err := fmt.Errorf("connection #%d: source output %s not found", i+1, conn.Source)

			errs = append(errs, err)
		}

		if _, found := blks.FindInput(conn.Target); !found {
			err := fmt.Errorf("connection #%d: target input %s not found", i+1, conn.Target)

			errs = append(errs, err)
		}
	}

	return errs
}
