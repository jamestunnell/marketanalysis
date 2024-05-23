package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/recorders"
)

func (a *Graphs) RunDay(w http.ResponseWriter, r *http.Request) {
	var runReq models.RunDayRequest
	if err := json.NewDecoder(r.Body).Decode(&runReq); err != nil {
		handleAppErr(w, app.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Debug().Interface("request", runReq).Msg("received run-day request")

	keyVal := mux.Vars(r)[a.Store.GetInfo().KeyName]

	cfg, appErr := a.Store.Get(r.Context(), keyVal)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	barsLoader := bars.NewAlpacaLoader(runReq.Symbol)
	buf := bytes.NewBuffer([]byte{})

	var rec blocks.Recorder
	var contentType string

	switch runReq.Format {
	case "csv":
		rec = recorders.NewCSV(buf, runReq.LocalTZ)
		contentType = "text/csv"
	case "json":
		rec = recorders.NewJSON(buf, runReq.LocalTZ)
		contentType = "application/json"
	default:
		appErr := app.NewErrInvalidInput("request", "format is missing or unknown")

		handleAppErr(w, appErr)

		return
	}

	err := graph.RunDay(runReq.Date, cfg, barsLoader, rec)
	if err != nil {
		appErr := app.NewErrActionFailed("run graph", err.Error())

		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(buf.Bytes()); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
