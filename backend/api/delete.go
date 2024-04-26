package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
)

func (a *API[T]) Delete(w http.ResponseWriter, r *http.Request) {
	keyVal := mux.Vars(r)[a.KeyName]

	result, err := a.Collection.DeleteOne(r.Context(), bson.D{{"_id", keyVal}})
	if err != nil {
		err = fmt.Errorf("failed to delete '%s': %w", keyVal, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	if result.DeletedCount == 0 {
		err = fmt.Errorf("document with key '%s' not found", keyVal)

		handleErr(w, err, http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
