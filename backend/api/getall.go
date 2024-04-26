package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (a *API[T]) GetAll(w http.ResponseWriter, r *http.Request) {
	cursor, err := a.Collection.Find(r.Context(), bson.D{})
	if err != nil {
		err = fmt.Errorf("failed to find %s: %w", a.NamePlural, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	var all []T

	err = cursor.All(r.Context(), &all)
	if err != nil {
		err = fmt.Errorf("failed to decode find results as %s: %w", a.NamePlural, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	// for no results
	if all == nil {
		all = []T{}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	p := map[string][]T{a.NamePlural: all}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		err = fmt.Errorf("failed to marshal response JSON: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}
}
