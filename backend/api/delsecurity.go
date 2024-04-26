package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type delSecurity struct {
	coll *mongo.Collection
}

func NewDelSecurity(coll *mongo.Collection) http.Handler {
	return &delSecurity{coll: coll}
}

func (h *delSecurity) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	result, err := h.coll.DeleteOne(r.Context(), bson.D{{"_id", symbol}})
	if err != nil {
		err = fmt.Errorf("failed to delete security '%s': %w", symbol, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	if result.DeletedCount == 0 {
		err = fmt.Errorf("security with symbol '%s' not found", symbol)

		handleErr(w, err, http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
