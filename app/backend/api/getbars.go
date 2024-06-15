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

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/app/backend/stores"
)

type getBars struct {
	db *mongo.Database
}

func NewGetBars(db *mongo.Database) http.Handler {
	return &getBars{db: db}
}

func (h *getBars) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	urlVals := r.URL.Query()
	dateStr := mux.Vars(r)["date"]
	timeZone := urlVals.Get("timeZone")

	d, err := date.Parse(date.RFC3339, dateStr)
	if err != nil {
		msg := fmt.Sprintf("date '%s'", dateStr)

		handleAppErr(w, backend.NewErrInvalidInput(msg, err.Error()))
	}

	var loc *time.Location

	if timeZone != "" {
		loc, err = time.LoadLocation(timeZone)
		if err != nil {
			msg := fmt.Sprintf("time zone '%s'", timeZone)

			handleAppErr(w, backend.NewErrInvalidInput(msg, err.Error()))
		}
	}

	store := stores.NewBarSets(symbol, h.db)

	barSet, appErr := store.Get(r.Context(), d.String())
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	if loc != nil {
		for _, bar := range barSet.Bars {
			bar.Timestamp = bar.Timestamp.In(loc)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(barSet); err != nil {
		log.Warn().Msg("failed to write response")
	}
}
