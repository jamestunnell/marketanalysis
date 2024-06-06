package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
)

func (a *Graphs) BacktestGraph(w http.ResponseWriter, r *http.Request) {
	keyVal := mux.Vars(r)[a.Store.GetInfo().KeyName]

	cfg, appErr := a.Store.Get(r.Context(), keyVal)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	var bt models.BacktestRequest
	if err := json.NewDecoder(r.Body).Decode(&bt); err != nil {
		handleAppErr(w, backend.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Debug().Interface("request", bt).Msg("received backtest request")

	loc, err := time.LoadLocation(bt.TimeZone)
	if err != nil {
		msg := fmt.Sprintf("time zone '%s'", bt.TimeZone)

		handleAppErr(w, backend.NewErrInvalidInput(msg, err.Error()))

		return
	}

	loader := backend.NewBarSetLoader(a.DB, bt.Symbol, loc)

	recording, err := graph.Backtest(
		r.Context(), cfg, bt.Symbol, bt.Date, loc, loader.Load, bt.ShowWarmup, bt.Predictor, bt.Threshold)
	if err != nil {
		appErr := backend.NewErrActionFailed("backtest graph", err.Error())

		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(recording); err != nil {
		log.Warn().Err(err).Msg("failed to encode response")
	}
}
