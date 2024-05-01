package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/recorders"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"
)

const (
	ParamNameSymbol = "symbol"
	ParamNameDate   = "date"
)

func (a *Graphs) Run(w http.ResponseWriter, r *http.Request) {
	urlVals := r.URL.Query()

	dateStr := urlVals.Get(ParamNameDate)
	if dateStr == "" {
		err := fmt.Errorf("date param is missing")

		handleErr(w, err, http.StatusBadRequest)

		return
	}

	runDate, err := date.Parse(date.RFC3339, dateStr)
	if err != nil {
		err := fmt.Errorf("run date %s is invalid: %w", dateStr, err)

		handleErr(w, err, http.StatusBadRequest)

		return
	}

	symbol := urlVals.Get(ParamNameSymbol)
	if symbol == "" {
		err := fmt.Errorf("symbol param is missing")

		handleErr(w, err, http.StatusBadRequest)

		return
	}

	security, herr := a.securities.FindOne(r.Context(), symbol)
	if herr != nil {
		handleErr(w, herr.Error, herr.StatusCode)

		return
	}

	loc, err := time.LoadLocation(security.TimeZone)
	if err != nil {
		err = fmt.Errorf("failed to load location from time zone '%s': %w", security.TimeZone, err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	keyVal := mux.Vars(r)[a.KeyName]

	cfg, herr := a.FindOne(r.Context(), keyVal)
	if herr != nil {
		handleErr(w, herr.Error, herr.StatusCode)

		return
	}

	buf := bytes.NewBuffer([]byte{})
	recorder := recorders.NewCSV(buf, loc)
	g := graph.New(cfg)

	if err = g.Init(recorder); err != nil {
		err = fmt.Errorf("failed to init graph: %w", err)

		handleErr(w, err, http.StatusBadRequest)

		return
	}

	barsLoader := bars.NewAlpacaLoader(security)
	if err = barsLoader.Init(); err != nil {
		err = fmt.Errorf("failed to init alpaca bars loader: %w", err)

		handleErr(w, err, http.StatusBadRequest)

		return
	}

	bars, err := barsLoader.GetRunBars(runDate, g.GetWarmupPeriod())
	if err != nil {
		err = fmt.Errorf("failed to get run bars: %w", err)

		handleErr(w, err, http.StatusInternalServerError)

		return
	}

	log.Debug().
		Stringer("date", runDate).
		Time("firstBar", bars[0].Timestamp).
		Time("lastBar", bars[len(bars)-1].Timestamp).
		Msgf("running model with %d bars", len(bars))

	for _, bar := range bars {
		g.Update(bar)
	}

	recorder.Flush()

	w.Header().Set("Content-Type", "text/csv")

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(buf.Bytes()); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}

func parseDateParam(urlVals url.Values) (date.Date, error) {
	dateStr := urlVals.Get(ParamNameDate)
	if dateStr == "" {
		return date.Date{}, fmt.Errorf("%s param is missing", ParamNameDate)
	}

	return date.Parse(date.RFC3339, dateStr)
}
