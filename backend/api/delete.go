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
		err = fmt.Errorf("failed to delete %s with %s '%s': %w", a.Name, a.KeyName, keyVal, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	if result.DeletedCount == 0 {
		err = fmt.Errorf("%s with %s '%s' not found", a.Name, a.KeyName, keyVal)

		handleErr(w, err, http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
