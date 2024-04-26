package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a *API[T]) Get(w http.ResponseWriter, r *http.Request) {
	keyVal := mux.Vars(r)[a.KeyName]

	var val T

	err := a.Collection.FindOne(r.Context(), bson.D{{"_id", keyVal}}).Decode(&val)
	if err == mongo.ErrNoDocuments {
		err = fmt.Errorf("%s with %s '%s' not found", a.Name, a.KeyName, keyVal)

		handleErr(w, err, http.StatusNotFound)

		return
	} else if err != nil {
		err = fmt.Errorf("failed to find %s with %s '%s': %w", a.Name, a.KeyName, keyVal, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(val); err != nil {
		log.Warn().Msg("failed to write response")
	}
}
