package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAll[T any](
	w http.ResponseWriter,
	r *http.Request,
	res *Resource[T],
	coll *mongo.Collection,
) {
	cursor, err := coll.Find(r.Context(), bson.D{})
	if err != nil {
		err = fmt.Errorf("failed to find %s: %w", res.NamePlural, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	var all []T

	err = cursor.All(r.Context(), &all)
	if err != nil {
		err = fmt.Errorf("failed to decode find results as %s: %w", res.NamePlural, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	// for no results
	if all == nil {
		all = []T{}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	p := map[string][]T{res.NamePlural: all}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		err = fmt.Errorf("failed to marshal response JSON: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}
}
