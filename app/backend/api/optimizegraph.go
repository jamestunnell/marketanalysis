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
)

func (a *Graphs) Optimize(w http.ResponseWriter, r *http.Request) {
	var opt bemodels.OptimizeRequest

	if err := json.NewDecoder(r.Body).Decode(&opt); err != nil {
		handleAppErr(w, backend.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Info().Interface("request", opt).Msg("received optimize request")

	job := &OptimizeJob{DB: a.DB, Request: &opt}

	if !a.BG.RunJob(job) {
		handleAppErr(w, backend.NewErrInvalidInput("job ID", "ID is already in use"))

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type OptimizeJob struct {
	DB      *mongo.Database
	Request *bemodels.OptimizeRequest
}

type OptimizeResponse struct {
}

func (job *OptimizeJob) GetID() string {
	return job.Request.JobID
}

func (job *OptimizeJob) Execute(onProgress background.JobProgressFunc) (any, error) {
	loader := backend.NewBarSetLoader(job.DB, job.Request.Symbol)
	log.Info().Msg("optimize job: started job")

	iter := 0
	maxIter := job.Request.OptimizeSettings.MaxIterations
	postEval := func(paramVals map[string]any, result float64) {
		iter++

		progress := float64(iter) / float64(maxIter)

		if iter%10 == 0 {
			log.Debug().Str("id", job.GetID()).Msgf("optimize job: %6.2f%%", progress*100.0)
		}

		onProgress(progress)
	}

	results, err := graph.Optimize(
		context.Background(),
		job.Request.Graph,
		job.Request.Days,
		job.Request.SourceQuantity,
		job.Request.TargetParams,
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
