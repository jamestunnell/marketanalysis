package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app/backend"
	bemodels "github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
)

const DefaultTimeZone = "America/New_York"

var defaultLoc *time.Location

func init() {
	var err error

	defaultLoc, err = time.LoadLocation(DefaultTimeZone)
	if err != nil {
		log.Warn().Err(err).Msg("failed to load default location")
	}
}

func (a *Graphs) Optimize(w http.ResponseWriter, r *http.Request) {
	var opt bemodels.OptimizeRequest

	if err := json.NewDecoder(r.Body).Decode(&opt); err != nil {
		handleAppErr(w, backend.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Debug().Interface("request", opt).Msg("received optimize request")

	loader := backend.NewBarSetLoader(a.DB, opt.Symbol)

	results, err := graph.Optimize(
		r.Context(), opt.Graph, opt.Days, opt.SourceQuantity,
		opt.TargetParams, opt.OptimizeSettings, loader.Load)
	if err != nil {
		appErr := backend.NewErrActionFailed("optimize graph", err.Error())

		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(results); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
