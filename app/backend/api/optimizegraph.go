package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/app/backend/background"
	bemodels "github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
)

const DefaultTimeZone = "America/New_York"

var defaultLoc *time.Location

func init() {
	var err error

	defaultLoc, err = time.LoadLocation(DefaultTimeZone)
	if err != nil {
		log.Warn().Err(err).Msg("failed to load default location")
	}
}

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

	// w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusAccepted)

	// if err = json.NewEncoder(w).Encode(results); err != nil {
	// 	log.Warn().Err(err).Msg("failed to write response")
	// }
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

func (job *OptimizeJob) Execute(progress background.JobProgressFunc) (any, error) {
	loader := backend.NewBarSetLoader(job.DB, job.Request.Symbol)
	log.Info().Msg("optimize: started job")

	results, err := graph.Optimize(
		context.Background(),
		job.Request.Graph,
		job.Request.Days,
		job.Request.SourceQuantity,
		job.Request.TargetParams,
		job.Request.OptimizeSettings,
		loader.Load,
	)
	if err != nil {
		log.Error().Err(err).Msg("optimize: failed")

		return nil, err
	}

	log.Info().Interface("results", results).Msg("optimize: job complete")

	return results, nil
}
