package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/rs/zerolog/log"
)

func Get[T backend.Resource](
	w http.ResponseWriter,
	r *http.Request,
	s backend.Store[T],
) {
	keyVal := mux.Vars(r)[s.GetInfo().KeyName]

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
