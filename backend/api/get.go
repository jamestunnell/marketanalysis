package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

func Get[T any](
	w http.ResponseWriter,
	r *http.Request,
	res *Resource[T],
	col *mongo.Collection,
) {
	keyVal := mux.Vars(r)[res.KeyName]

	val, herr := FindOne[T](r.Context(), keyVal, res, col)
	if herr != nil {
		handleErr(w, herr.Error, herr.StatusCode)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(val); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
