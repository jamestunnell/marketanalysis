package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-multierror"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/xeipuuv/gojsonschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *API[T]) Put(w http.ResponseWriter, r *http.Request) {
	urlKeyVal := mux.Vars(r)[a.KeyName]

	d, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read request body: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	// the key value can be left out of the request JSON since it's in the URL
	requestKeyVal := gjson.GetBytes(d, a.KeyName).String()
	if requestKeyVal == "" {
		d, err = sjson.SetBytes(d, a.KeyName, urlKeyVal)
		if err != nil {
			err = fmt.Errorf("failed to insert %s '%s' into JSON: %w", a.KeyName, urlKeyVal, err)

			handleErr(w, err, http.StatusInternalServerError)

			return
		}
	} else if requestKeyVal != urlKeyVal {
		err = fmt.Errorf("%s '%s' in JSON does not match '%s' in URL", a.KeyName, requestKeyVal, urlKeyVal)

		handleErr(w, err, http.StatusBadRequest)

		return
	}

	vResult, err := a.Schema.Validate(gojsonschema.NewBytesLoader(d))
	if err != nil {
		err = fmt.Errorf("failed to validate request JSON: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	if !vResult.Valid() {
		var merr *multierror.Error

		for _, resultErr := range vResult.Errors() {
			merr = multierror.Append(merr, fmt.Errorf("%s", resultErr.String()))
		}

		err = fmt.Errorf("%s JSON is invalid: %w", a.Name, merr)

		handleErr(w, err, http.StatusBadRequest)

		return
	}

	var val T

	err = json.Unmarshal(d, &val)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal request JSON: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	if err = a.Validate(&val); err != nil {
		err = fmt.Errorf("unmarshaled %s is invalid: %w", a.Name, err)

		handleErr(w, err, http.StatusBadRequest)

		return
	}

	_, err = a.Collection.ReplaceOne(
		r.Context(), bson.D{{"_id", urlKeyVal}}, val, options.Replace().SetUpsert(true))
	if err != nil {
		err = fmt.Errorf("failed to upsert %s into collection: %w", a.Name, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
