package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/xeipuuv/gojsonschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jamestunnell/marketanalysis/backend/models"
)

type putSecurity struct {
	coll *mongo.Collection
}

func NewPutSecurity(coll *mongo.Collection) http.Handler {
	return &putSecurity{coll: coll}
}

func (h *putSecurity) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]

	d, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("failed to read request body: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	// the symbol can be left out of the request JSON since it's in the URL
	reqSymbol := gjson.GetBytes(d, "symbol").String()

	if reqSymbol == "" {
		log.Debug().Str("symbol", symbol).Msg("inserting symbol into JSON")
		d, err = sjson.SetBytes(d, "symbol", symbol)
		if err != nil {
			err = fmt.Errorf("failed to insert symbol into JSON: %w", err)

			handleErr(w, err, http.StatusInternalServerError)

			return
		}
	} else if reqSymbol != symbol {
		err = fmt.Errorf("symbol '%s' in request JSON does not match symbol '%s' in URL", reqSymbol, symbol)

		handleErr(w, err, http.StatusBadRequest)

		return
	}

	var security models.Security

	if err = json.Unmarshal(d, &security); err != nil {
		err = fmt.Errorf("failed to unmarshal request JSON: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	schema, err := models.LoadSecuritySchema()
	if err != nil {
		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	vResult, err := schema.Validate(gojsonschema.NewBytesLoader(d))
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

		handleErr(w, merr, http.StatusBadRequest)

		return
	}

	_, err = h.coll.ReplaceOne(
		r.Context(), bson.D{{"_id", symbol}}, security, options.Replace().SetUpsert(true))
	if err != nil {
		err = fmt.Errorf("failed to upsert security into collection: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
