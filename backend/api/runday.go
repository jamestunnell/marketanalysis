package api

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/recorders"
)

type RunDayRequest struct {
	Symbol string    `json:"symbol"`
	Date   date.Date `json:"date"`
}

func (a *Graphs) RunDay(w http.ResponseWriter, r *http.Request) {
	runReq, err := LoadRequestJSON[RunDayRequest](r)
	if err != nil {
		appErr := app.NewErrActionFailed("load request JSON", err.Error())

		handleAppErr(w, appErr)

		return
	}

	security, appErr := a.securities.Store.Get(r.Context(), runReq.Symbol)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	loc, err := time.LoadLocation(security.TimeZone)
	if err != nil {
		action := fmt.Sprintf("load location from time zone '%s'", security.TimeZone)

		handleAppErr(w, app.NewErrActionFailed(action, err.Error()))

		return
	}

	keyVal := mux.Vars(r)[a.Store.RDef().KeyName]

	cfg, appErr := a.Store.Get(r.Context(), keyVal)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	barsLoader := bars.NewAlpacaLoader(security)
	buf := bytes.NewBuffer([]byte{})
	recorder := recorders.NewCSV(buf, loc)

	err = graph.RunDay(security, runReq.Date, cfg, barsLoader, recorder)
	if err != nil {
		appErr := app.NewErrActionFailed("run graph", err.Error())

		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", "text/csv")

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(buf.Bytes()); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
