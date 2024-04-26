package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamestunnell/marketanalysis/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type getSecurities struct {
	coll *mongo.Collection
}

type GetSecuritiesResponse struct {
	Securities []models.Security `json:"securities"`
}

func NewGetSecurities(coll *mongo.Collection) http.Handler {
	return &getSecurities{coll: coll}
}

func (h *getSecurities) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cursor, err := h.coll.Find(r.Context(), bson.D{})
	if err != nil {
		err = fmt.Errorf("failed to find securities: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	var securities []models.Security

	err = cursor.All(r.Context(), &securities)
	if err != nil {
		err = fmt.Errorf("failed to decode find results: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	p := GetSecuritiesResponse{Securities: securities}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		err = fmt.Errorf("failed to marshal response JSON: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}
}
