package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type getSecurity struct {
	coll *mongo.Collection
}

func NewGetSecurity(coll *mongo.Collection) http.Handler {
	return &getSecurity{coll: coll}
}

func (h *getSecurity) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	var security models.Security

	err := h.coll.FindOne(r.Context(), bson.D{{"_id", symbol}}).Decode(&security)
	if err == mongo.ErrNoDocuments {
		err = fmt.Errorf("security with symbol '%s' not found", symbol)

		handleErr(w, err, http.StatusNotFound)

		return
	} else if err != nil {
		err = fmt.Errorf("failed to find security: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(security); err != nil {
		log.Warn().Msg("failed to write response")
	}
}
