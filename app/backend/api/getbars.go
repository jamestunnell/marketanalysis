package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app"
)

type getBars struct {
	db *mongo.Database
}

const DefaultTimeZone = "America/New_York"

func NewGetBars(db *mongo.Database) http.Handler {
	return &getBars{db: db}
}

func (h *getBars) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	urlVals := r.URL.Query()
	timeZone := urlVals.Get("timeZone")

	if timeZone == "" {
		timeZone = DefaultTimeZone
	}

	d, err := date.Parse(date.RFC3339, mux.Vars(r)["date"])

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		msg := fmt.Sprintf("time zone '%s'", timeZone)

		handleAppErr(w, app.NewErrInvalidInput(msg, err.Error()))

		return
	}

	loader := app.NewDayBarsLoader(h.db, symbol, loc)

	dayBars, err := loader.Load(r.Context(), d)
	if err != nil {
		handleAppErr(w, app.NewErrActionFailed("load day bars", err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(dayBars); err != nil {
		log.Warn().Msg("failed to write response")
	}
}
