package backend

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/backend/api"
)

func BindAPI(r *mux.Router, db *mongo.Database) {
	securities, err := api.NewSecurities(db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make securities API")
	}

	graphs, err := api.NewGraphs(db, securities)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make securities API")
	}

	r.Handle("/blocks/{type}", api.NewGetBlockInfo()).Methods(http.MethodGet)
	r.Handle("/blocks", api.NewGetAllBlockInfo()).Methods(http.MethodGet)
	r.Handle("/status", api.NewStatus())

	graphs.Bind(r)
	securities.Bind(r)
}
