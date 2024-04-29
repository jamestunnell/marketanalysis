package backend

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/backend/api"
)

func BindAPI(r *mux.Router, db *mongo.Database) {
	r.Handle("/status", api.NewStatus())

	r.Handle("/blocks/{type}", api.NewGetBlockInfo()).Methods(http.MethodGet)
	r.Handle("/blocks", api.NewGetAllBlockInfo()).Methods(http.MethodGet)

	securitiesAPI, err := api.MakeSecuritiesAPI(db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make securities API")
	}

	graphsAPI, err := api.MakeGraphsAPI(db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make securities API")
	}

	securitiesAPI.Bind(r)
	graphsAPI.Bind(r)
}
