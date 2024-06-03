package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/graph"
)

func (a *Graphs) EvalGraph(w http.ResponseWriter, r *http.Request) {
	keyVal := mux.Vars(r)[a.Store.GetInfo().KeyName]

	cfg, appErr := a.Store.Get(r.Context(), keyVal)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		handleAppErr(w, app.NewErrInvalidInput("eval request body", err.Error()))
	}

	runType := gjson.GetBytes(d, "type")
	if !runType.Exists() {
		handleAppErr(w, app.NewErrInvalidInput("eval request JSON", "missing type"))
	}

	switch runType.String() {
	case models.EvalSlope:
		a.EvalSlope(r.Context(), w, cfg, d)
	}
}

func (a *Graphs) EvalSlope(
	ctx context.Context,
	w http.ResponseWriter,
	cfg *graph.Configuration,
	requestData []byte,
) {
	var eval models.EvalSlopeRequest
	if err := json.Unmarshal(requestData, &eval); err != nil {
		handleAppErr(w, app.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	log.Debug().Interface("request", eval).Msg("received eval-slope request")

	loc, err := time.LoadLocation(eval.TimeZone)
	if err != nil {
		handleAppErr(w, app.NewErrInvalidInput("eval timeZone", err.Error()))

		return
	}

	loader := app.NewDayBarsLoader(a.DB, eval.Symbol, loc)

	recording, err := graph.EvalSlope(
		ctx, cfg, eval.Symbol, eval.Date,
		loc, loader, eval.Source, eval.Predictor, eval.Horizon)
	if err != nil {
		appErr := app.NewErrActionFailed("eval graph", err.Error())

		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(recording); err != nil {
		log.Warn().Err(err).Msg("failed to encode response")
	}
}
