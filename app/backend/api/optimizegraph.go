package api

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/maps"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/app/backend/background"
	bemodels "github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
	"github.com/jamestunnell/marketanalysis/models"
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

type OptimizeResult struct {
	Score     float64
	ParamVals models.ParamVals
}

func (job *OptimizeGraphParamsJob) GetID() string {
	return job.Request.JobID
}

func (job *OptimizeGraphParamsJob) Execute(onProgress background.JobProgressFunc) (any, error) {
	log.Info().Msg("optimize job: started job")

	iter := 0
	maxIter := job.Request.OptimizeSettings.MaxIterations
	postEval := func(result *optimization.Result[models.ParamVals]) {
		iter++

		progress := float64(iter) / float64(maxIter)

		log.Debug().
			Str("id", job.GetID()).
			Float64("score", result.Score).
			Interface("paramVals", result.Value).
			Str("progress", fmt.Sprintf("%6.2f%%", progress*100.0)).
			Msg("optimize job: progress update")

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

	log.Info().Interface("best", results.Best).Msg("optimize job: complete")

	job.writeReport(results)

	return results, nil
}

func (job *OptimizeGraphParamsJob) writeReport(results *optimization.Results[models.ParamVals]) {
	reportName := fmt.Sprintf("%s-%s.csv", job.Request.Graph.Name, time.Now().Format(time.RFC3339))

	f, err := os.Create(reportName)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create report file")

		return
	}

	w := csv.NewWriter(f)

	defer f.Close()
	defer w.Flush()

	keys := maps.Keys(results.Last.Value)

	slices.Sort(keys)

	header := append([]string{"iteration", "score"}, keys...)
	iterIdx := 0
	scoreIdx := 1
	keysOffset := 2

	_ = w.Write(header)

	record := make([]string, len(header))

	for i, result := range results.History {
		record[iterIdx] = strconv.Itoa(i + 1)
		record[scoreIdx] = strconv.FormatFloat(result.Score, 'f', 3, 64)

		for keyIdx, key := range keys {
			record[keyIdx+keysOffset] = fmt.Sprintf("%v", result.Value[key])
		}

		_ = w.Write(record)
	}

	log.Info().
		Str("fname", f.Name()).
		Str("jobID", job.Request.JobID).
		Msg("wrote optimization report file")
}
