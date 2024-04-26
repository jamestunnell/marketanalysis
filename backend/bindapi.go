package backend

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/backend/api"
	"github.com/jamestunnell/marketanalysis/backend/models"
)

func BindAPI(r *mux.Router, db *mongo.Database) {
	r.Handle("/status", api.NewStatus())

	schema, err := models.LoadSchema(models.SecuritySchemaStr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load security schema")
	}

	securitiesAPI := &api.API[models.Security]{
		KeyName:    models.SecurityKeyName,
		NamePlural: models.SecurityNamePlural,
		Collection: db.Collection(models.SecurityNamePlural),
		Schema:     schema,
	}

	schema, err = models.LoadSchema(models.ModelDefSchemaStr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load model def schema")
	}

	modelDefsAPI := &api.API[models.ModelDef]{
		KeyName:    models.ModelDefKeyName,
		NamePlural: models.ModelDefNamePlural,
		Collection: db.Collection(models.ModelDefNamePlural),
		Schema:     schema,
	}

	securitiesAPI.Bind(r)
	modelDefsAPI.Bind(r)
}
