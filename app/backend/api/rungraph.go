package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app/backend"
	bemodels "github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/models"
)

func (a *Graphs) Run(w http.ResponseWriter, r *http.Request) {
	var run bemodels.RunRequest

	if err := json.NewDecoder(r.Body).Decode(&run); err != nil {
		handleAppErr(w, backend.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Debug().Interface("request", run).Msg("received run request")

	loader := backend.NewBarSetLoader(a.DB, run.Symbol)

	var timeSeries *models.TimeSeries

	var err error

	switch run.RunType {
	case bemodels.RunSingleDay:
		timeSeries, err = graph.RunSingleDay(
			r.Context(), run.Graph, run.Date, loader.Load)
	case bemodels.RunMultiDay:
		timeSeries, err = graph.RunMultiDay(
			r.Context(), run.Graph, run.Date, loader.Load)
	case bemodels.RunMultiDaySummary:
		timeSeries, err = graph.RunMultiDaySummary(
			r.Context(), run.Graph, run.Date, loader.Load)
	default:
		msg := fmt.Sprintf("run type '%s'", run.RunType)

		handleAppErr(w, backend.NewErrInvalidInput(msg))

		return
	}

	if err != nil {
		appErr := backend.NewErrActionFailed("run graph", err.Error())

		handleAppErr(w, appErr)

		return
	}

	// do clustering by mean and stddev
	if run.NumCharts > 0 && !timeSeries.IsEmpty() {
		err := timeSeries.Cluster(run.NumCharts, models.QuantityMeanStddev)
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
