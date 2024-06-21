package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/app/backend/background"
)

type JobStatus struct {
	bg background.System
}

func NewJobStatus(bg background.System) *JobStatus {
	return &JobStatus{bg: bg}
}

func (h *JobStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]

	s, found := h.bg.GetJobStatus(id)
	if !found {
		handleAppErr(w, backend.NewErrNotFound("job ID"))

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(s); err != nil {
		log.Warn().Msg("failed to write response")
	}
}
