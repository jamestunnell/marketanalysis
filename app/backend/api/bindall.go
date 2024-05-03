package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

func BindAll(r *mux.Router, db *mongo.Database) {
	securities, err := NewSecurities(db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make securities API")
	}

	graphs, err := NewGraphs(db, securities)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make securities API")
	}

	r.Handle("/blocks/{type}", NewGetBlockInfo()).Methods(http.MethodGet)
	r.Handle("/blocks", NewGetAllBlockInfo()).Methods(http.MethodGet)
	r.Handle("/status", NewStatus())

	graphs.Bind(r)
	securities.Bind(r)
}
