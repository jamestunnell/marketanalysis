package api

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/backend/models"
	"github.com/jamestunnell/marketanalysis/blocks/registry"
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

				if _, found := registry.Get(b.Type); !found {
					return fmt.Errorf("unknown block type '%s'", b.Type)
				}
			}

			return nil
		},
	}

	return a, nil
}
