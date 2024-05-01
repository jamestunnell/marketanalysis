package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app"
)

func Get[T any](
	w http.ResponseWriter,
	r *http.Request,
	s app.Store[T],
) {
	keyVal := mux.Vars(r)[s.RDef().KeyName]

	val, appErr := s.Get(r.Context(), keyVal)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(val); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
