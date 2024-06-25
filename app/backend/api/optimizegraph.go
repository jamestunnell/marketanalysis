package api

import (
	"context"
	"encoding/json"
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

	job := &OptimizeGraphParamsJob{DB: a.DB, Request: &opt}

	if !a.BG.RunJob(job) {
		handleAppErr(w, backend.NewErrInvalidInput("job ID", "ID is already in use"))

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type OptimizeGraphParamsJob struct {
	DB      *mongo.Database
	Request *bemodels.OptimizeGraphParamsRequest
}

type OptimizeResponse struct {
}

func (job *OptimizeGraphParamsJob) GetID() string {
	return job.Request.JobID
}

func (job *OptimizeGraphParamsJob) Execute(onProgress background.JobProgressFunc) (any, error) {
	loader := backend.NewBarSetLoader(job.DB, job.Request.Symbol)
	log.Info().Msg("optimize job: started job")

	iter := 0
	maxIter := job.Request.OptimizeSettings.MaxIterations
	postEval := func(result *optimization.Result) {
		iter++

		progress := float64(iter) / float64(maxIter)

		if iter%10 == 0 {
			log.Debug().Str("id", job.GetID()).Msgf("optimize job: %6.2f%%", progress*100.0)
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
		loader.Load,
		postEval,
	)
	if err != nil {
		log.Error().Err(err).Msg("optimize job: failed")

		return nil, err
	}

	log.Info().Interface("results", results).Msg("optimize job: complete")

	return results, nil
}
