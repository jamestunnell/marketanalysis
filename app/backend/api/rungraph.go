package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/recorders"
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
		handleAppErr(w, app.NewErrInvalidInput("run request body", err.Error()))
	}

	runType := gjson.GetBytes(d, "type")
	if !runType.Exists() {
		handleAppErr(w, app.NewErrInvalidInput("run request JSON", "missing type"))
	}

	switch runType.String() {
	case models.RunDay:
		a.RunDay(w, cfg, d)
	}
}

func (a *Graphs) RunDay(
	w http.ResponseWriter,
	cfg *graph.Configuration,
	requestData []byte,
) {
	var runDay models.RunDayRequest
	if err := json.Unmarshal(requestData, &runDay); err != nil {
		handleAppErr(w, app.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Debug().Interface("request", runDay).Msg("received run-day request")

	barsLoader := bars.NewAlpacaLoader(runDay.Symbol)
	buf := bytes.NewBuffer([]byte{})

	var rec blocks.Recorder
	var contentType string

	switch runDay.Format {
	case "csv":
		rec = recorders.NewCSV(buf, runDay.LocalTZ)
		contentType = "text/csv"
	case "json":
		rec = recorders.NewTimeSeriesJSON(buf, runDay.LocalTZ)
		contentType = "application/json"
	default:
		appErr := app.NewErrInvalidInput("request", "format is missing or unknown")

		handleAppErr(w, appErr)

		return
	}

	err := graph.RunDay(runDay.Date, cfg, barsLoader, rec)
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
