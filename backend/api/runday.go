package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-multierror"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/recorders"
)

const (
	ParamNameSymbol = "symbol"
	ParamNameDate   = "date"
)

func (a *Graphs) RunDay(w http.ResponseWriter, r *http.Request) {
	symbol, runDate, err := parseRunParams(r.URL.Query())
	if err != nil {
		err := fmt.Errorf("invalid query params: %w", err)

		handleAppErr(w, &app.Error{Err: err, Code: app.InvalidInput})

		return
	}

	security, appErr := a.securities.Store.Get(r.Context(), symbol)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	loc, err := time.LoadLocation(security.TimeZone)
	if err != nil {
		action := fmt.Sprintf("load location from time zone '%s'", security.TimeZone)

		handleAppErr(w, app.NewActionFailedError(action, err))

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

	err = graph.RunDay(security, runDate, cfg, barsLoader, recorder)
	if err != nil {
		appErr := app.NewActionFailedError("run graph", err)

		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", "text/csv")

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(buf.Bytes()); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}

func parseRunParams(urlVals url.Values) (symbol string, runDate date.Date, err error) {
	symbol = urlVals.Get(ParamNameSymbol)
	dateStr := urlVals.Get(ParamNameDate)
	errs := []error{}

	if dateStr != "" {
		var parseErr error

		runDate, parseErr = date.Parse(date.RFC3339, dateStr)
		if parseErr != nil {
			errs = append(errs, fmt.Errorf("invalid date values '%s': %w", dateStr, parseErr))
		}
	} else {
		errs = append(errs, fmt.Errorf("%s param is missing", ParamNameDate))
	}

	if symbol == "" {
		errs = append(errs, fmt.Errorf("%s param is missing", ParamNameSymbol))
	}

	if len(errs) > 0 {
		var merr *multierror.Error

		for _, oneErr := range errs {
			merr = multierror.Append(merr, oneErr)
		}

		err = merr
	}

	return
}
