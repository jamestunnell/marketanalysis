package graph

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/xeipuuv/gojsonschema"

	"github.com/jamestunnell/marketanalysis/blocks/registry"
)

const ConfigKeyName = "id"

type Configuration struct {
	ID          string                  `json:"id" bson:"_id"`
	Name        string                  `json:"name"`
	Blocks      map[string]*BlockConfig `json:"blocks"`
	Connections []*Connection           `json:"connections"`
}

type BlockConfig struct {
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

	for blkName, b := range cfg.Blocks {
		newBlk, blkFound := registry.Get(b.Type)
		if !blkFound {
			err := fmt.Errorf("block %s: has unknown type '%s'", blkName, b.Type)

			errs = append(errs, err)

			continue
		}

		blk := newBlk()

		blks[blkName] = blk

		// TODO: validate recording outputs

		// validate params
		params := blk.GetParams()

		for name, val := range b.ParamVals {
			param, found := params[name]
			if !found {
				err := fmt.Errorf("block %s: unknown param name '%s'", blkName, name)

				errs = append(errs, err)

				continue
			}

			l := gojsonschema.NewGoLoader(param.GetSchema())

			schema, err := gojsonschema.NewSchema(l)
			if err != nil {
				err := fmt.Errorf("block %s: failed to compile schema for param %s: %w", blkName, name, err)

				errs = append(errs, err)

				continue
			}

			result, err := schema.Validate(gojsonschema.NewGoLoader(val))
			if err != nil {
				err := fmt.Errorf("block %s: failed to validate value %v for param %s: %w", blkName, val, name, err)

				errs = append(errs, err)

				continue
			}

			if !result.Valid() {
				errs = append(errs, newValidateParamValErr(name, val, result))
			}
		}
	}

	// TODO: validate connections

	return errs
}

func newValidateParamValErr(
	name string,
	val any,
	result *gojsonschema.Result,
) error {
	var merr *multierror.Error

	for _, resultErr := range result.Errors() {
		merr = multierror.Append(merr, fmt.Errorf("%s", resultErr.String()))
	}

	return fmt.Errorf("param %s value %v is invalid: %w", name, val, merr)
}
