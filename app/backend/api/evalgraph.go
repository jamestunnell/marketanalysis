package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/bars"
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
		a.EvalSlope(w, cfg, d)
	}
}

func (a *Graphs) EvalSlope(
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

	barsLoader := bars.NewAlpacaLoader(eval.Symbol)

	recording, err := graph.EvalSlope(cfg, barsLoader, eval.EvalSlopeConfig)
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
