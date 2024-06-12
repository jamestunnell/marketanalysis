package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app/backend"
	bemodels "github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/models"
)

func (a *Graphs) RunDay(w http.ResponseWriter, r *http.Request) {
	var runDay bemodels.RunDayRequest

	if err := json.NewDecoder(r.Body).Decode(&runDay); err != nil {
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
		r.Context(), runDay.Graph, runDay.Symbol, runDay.Date, loc, loader.Load)
	if err != nil {
		appErr := backend.NewErrActionFailed("run graph", err.Error())

		handleAppErr(w, appErr)

		return
	}

	// do clustering by mean and stddev
	if runDay.NumCharts > 0 && !timeSeries.IsEmpty() {
		err := timeSeries.Cluster(runDay.NumCharts, models.QuantityMeanStddev)
		if err != nil {
			appErr := backend.NewErrActionFailed("cluster results", err.Error())

			handleAppErr(w, appErr)

			return
		}

		clusters := map[string]int{}

		for _, q := range timeSeries.Quantities {
			cluster, found := q.Attributes[models.AttrCluster]
			if !found {
				continue
			}

			clusters[q.Name] = cluster.(int)
		}

		log.Debug().Interface("clusters", clusters).Msg("quantity clusters")
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(timeSeries); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
