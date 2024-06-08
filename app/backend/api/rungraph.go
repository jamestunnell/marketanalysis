package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app/backend"
	bemodels "github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
)

func (a *Graphs) RunDay(w http.ResponseWriter, r *http.Request) {
	var runDay bemodels.RunDayRequest

	if err := json.NewDecoder(r.Body).Decode(&runDay); err != nil {
		handleAppErr(w, backend.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Debug().Interface("request", runDay).Msg("received run-day request")

	loc, err := time.LoadLocation(runDay.TimeZone)
	if err != nil {
		msg := fmt.Sprintf("run time zone '%s'", runDay.TimeZone)

		handleAppErr(w, backend.NewErrInvalidInput(msg, err.Error()))

		return
	}

	loader := backend.NewBarSetLoader(a.DB, runDay.Symbol, loc)

	timeSeries, err := graph.RunDay(
		r.Context(), runDay.Graph, runDay.Symbol, runDay.Date, loc, loader.Load)
	if err != nil {
		appErr := backend.NewErrActionFailed("run graph", err.Error())

		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(timeSeries); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
