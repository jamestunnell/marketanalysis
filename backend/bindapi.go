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

	schema, err := models.LoadSecuritySchema()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load security schema")
	}

	securitiesAPI := &api.API[models.Security]{
		KeyName:    "symbol",
		NamePlural: "securities",
		Collection: db.Collection("securities"),
		Schema:     schema,
	}

	securitiesAPI.Bind(r)
}
