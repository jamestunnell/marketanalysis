package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/bars"
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
		handleAppErr(w, app.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Debug().Interface("request", bt).Msg("received backtest request")

	loc, err := time.LoadLocation(bt.TimeZone)
	if err != nil {
		handleAppErr(w, app.NewErrInvalidInput("backtest timeZone", err.Error()))

		return
	}

	recording, err := graph.Backtest(cfg, bt.Symbol, bt.Date, loc, bars.GetAlpacaBarsOneMin, bt.Predictor, bt.Threshold)
	if err != nil {
		appErr := app.NewErrActionFailed("backtest graph", err.Error())

		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(recording); err != nil {
		log.Warn().Err(err).Msg("failed to encode response")
	}
}
