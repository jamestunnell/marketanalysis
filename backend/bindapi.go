package backend

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/backend/api"
	"github.com/jamestunnell/marketanalysis/backend/models"
	"github.com/jamestunnell/marketanalysis/blocks"
)

func BindAPI(r *mux.Router, db *mongo.Database) {
	r.Handle("/status", api.NewStatus())

	schema, err := models.LoadSchema(models.SecuritySchemaStr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load security schema")
	}

	securitiesAPI := &api.API[models.Security]{
		KeyName:    models.SecurityKeyName,
		Name:       models.SecurityName,
		NamePlural: models.SecurityNamePlural,
		Collection: db.Collection(models.SecurityNamePlural),
		Schema:     schema,
		Validate: func(s *models.Security) error {
			if s.Open.MinuteOfDay() >= s.Close.MinuteOfDay() {
				return fmt.Errorf("open '%s' is not before close '%s'", s.Open, s.Close)
			}

			return nil
		},
	}

	schema, err = models.LoadSchema(models.GraphSchemaStr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load model def schema")
	}

	grapsAPI := &api.API[models.Graph]{
		KeyName:    models.GraphKeyName,
		Name:       models.GraphName,
		NamePlural: models.GraphNamePlural,
		Collection: db.Collection(models.GraphNamePlural),
		Schema:     schema,
		Validate: func(g *models.Graph) error {
			names := map[string]int{}
			for _, b := range g.Blocks {
				if _, found := names[b.Name]; found {
					return fmt.Errorf("block name '%s' is not unique", b.Name)
				}

				if _, found := blocks.Registry().Get(b.Type); !found {
					return fmt.Errorf("unknown block type '%s'", b.Type)
				}
			}

			return nil
		},
	}

	securitiesAPI.Bind(r)
	grapsAPI.Bind(r)
}
