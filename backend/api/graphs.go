package api

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hashicorp/go-multierror"
	"github.com/jamestunnell/marketanalysis/backend/models"
	"github.com/jamestunnell/marketanalysis/blocks/registry"
	"github.com/xeipuuv/gojsonschema"
)

func MakeGraphsAPI(db *mongo.Database) (*API[models.GraphDef], error) {
	schema, err := models.LoadSchema(models.GraphDefSchemaStr)
	if err != nil {
		return nil, fmt.Errorf("failed to load graph def schema: %w", err)
	}

	a := &API[models.GraphDef]{
		KeyName:    models.GraphDefKeyName,
		Name:       models.GraphDefName,
		NamePlural: models.GraphDefNamePlural,
		Collection: db.Collection(models.GraphDefNamePlural),
		Schema:     schema,
		Validate: func(g *models.GraphDef) error {
			names := map[string]int{}
			for _, b := range g.Blocks {
				if _, found := names[b.Name]; found {
					return fmt.Errorf("block name '%s' is not unique", b.Name)
				}

				names[b.Name] = 1

				newBlk, blkFound := registry.Get(b.Type)
				if !blkFound {
					return fmt.Errorf("unknown block type '%s'", b.Type)
				}

				params := newBlk().GetParams()

				for name, val := range b.ParamVals {
					param, found := params[name]
					if !found {
						return fmt.Errorf("unknown param name '%s'", name)
					}

					l := gojsonschema.NewGoLoader(param.GetSchema())

					schema, err := gojsonschema.NewSchema(l)
					if err != nil {
						return fmt.Errorf("failed to compile schema for param %s: %w", name, err)
					}

					result, err := schema.Validate(gojsonschema.NewGoLoader(val))
					if err != nil {
						return fmt.Errorf("failed to validate value %v for param %s: %w", val, name, err)
					}

					if !result.Valid() {
						return newValidateParamValErr(name, val, result)
					}
				}
			}

			return nil
		},
	}

	return a, nil
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
