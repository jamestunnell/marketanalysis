package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	"github.com/jamestunnell/marketanalysis/app/backend"
	bemodels "github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
)

func (a *Graphs) RunGraph(w http.ResponseWriter, r *http.Request) {
	keyVal := mux.Vars(r)[a.Store.GetInfo().KeyName]

	cfg, appErr := a.Store.Get(r.Context(), keyVal)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		handleAppErr(w, backend.NewErrInvalidInput("run request body", err.Error()))
	}

	runType := gjson.GetBytes(d, "type")
	if !runType.Exists() {
		handleAppErr(w, backend.NewErrInvalidInput("run request JSON", "missing type"))
	}

	switch runType.String() {
	case bemodels.RunDay:
		a.RunDay(r.Context(), w, cfg, d)
	default:
		msg := fmt.Sprintf("run type '%s'", runType)

		handleAppErr(w, backend.NewErrInvalidInput(msg, "type is unknown"))
	}
}

func (a *Graphs) RunDay(
	ctx context.Context,
	w http.ResponseWriter,
	cfg *graph.Configuration,
	requestData []byte,
) {
	var runDay bemodels.RunDayRequest
	if err := json.Unmarshal(requestData, &runDay); err != nil {
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
		ctx, cfg, runDay.Symbol, runDay.Date, loc, loader.Load, runDay.ShowWarmup)
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
