package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/app/backend/background"
	bemodels "github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/optimization"
)

func (a *Graphs) OptimizeParams(w http.ResponseWriter, r *http.Request) {
	var opt bemodels.OptimizeGraphParamsRequest

	if err := json.NewDecoder(r.Body).Decode(&opt); err != nil {
		handleAppErr(w, backend.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Info().Interface("request", opt).Msg("received optimize request")

	loader, err := backend.NewBarSetLoader(a.DB, opt.Symbol)
	if err != nil {
		handleAppErr(w, backend.NewErrActionFailed("make bar loader", err.Error()))

		return
	}

	job := &OptimizeGraphParamsJob{DB: a.DB, Request: opt, Load: loader.Load}

	if !a.BG.RunJob(job) {
		handleAppErr(w, backend.NewErrInvalidInput("job ID", "ID is already in use"))

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type OptimizeGraphParamsJob struct {
	DB      *mongo.Database
	Load    graph.LoadBarsFunc
	Request bemodels.OptimizeGraphParamsRequest
}

type OptimizeResponse struct {
}

func (job *OptimizeGraphParamsJob) GetID() string {
	return job.Request.JobID
}

func (job *OptimizeGraphParamsJob) Execute(onProgress background.JobProgressFunc) (any, error) {
	log.Info().Msg("optimize job: started job")

	iter := 0
	maxIter := job.Request.OptimizeSettings.MaxIterations
	postEval := func(result *optimization.Result) {
		iter++

		progress := float64(iter) / float64(maxIter)

		if iter%10 == 0 {
			log.Debug().
				Str("id", job.GetID()).
				Interface("result", result).
				Str("progress", fmt.Sprintf("%6.2f%%", progress*100.0)).
				Msgf("optimize job: progress update")
		}

		onProgress(progress)
	}

	results, err := graph.OptimizeParameters(
		context.Background(),
		job.Request.Graph,
		job.Request.Days,
		job.Request.SourceQuantity,
		job.Request.TargetParams,
		job.Request.ObjectiveType,
		job.Request.OptimizeSettings,
		job.Load,
		postEval,
	)
	if err != nil {
		log.Error().Err(err).Msg("optimize job: failed")

		return nil, err
	}

	log.Info().Interface("final result", results.Result).Msg("optimize job: complete")

	return results, nil
}
