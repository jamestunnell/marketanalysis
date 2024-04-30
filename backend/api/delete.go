package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Delete[T any](
	w http.ResponseWriter,
	r *http.Request,
	res *Resource[T],
	col *mongo.Collection,
) {
	keyVal := mux.Vars(r)[res.KeyName]

	result, err := col.DeleteOne(r.Context(), bson.D{{"_id", keyVal}})
	if err != nil {
		err = fmt.Errorf("failed to delete %s with %s '%s': %w", res.Name, res.KeyName, keyVal, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	if result.DeletedCount == 0 {
		err = fmt.Errorf("%s with %s '%s' not found", res.Name, res.KeyName, keyVal)

		handleErr(w, err, http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
