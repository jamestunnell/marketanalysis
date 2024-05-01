package api

import (
	"encoding/json"
	"net/http"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	ErrType string   `json:"errType"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func handleAppErr(w http.ResponseWriter, appErr app.Error) {
	var statusCode int

	switch appErr.GetType() {
	case app.NotFound:
		statusCode = http.StatusNotFound
	case app.InvalidInput:
		statusCode = http.StatusBadRequest
	case app.ActionFailed:
		statusCode = http.StatusInternalServerError
	default:
		log.Error().Msgf("app error type %s is unknown", appErr.GetType())

		statusCode = http.StatusInternalServerError
	}

	resp := &ErrorResponse{
		ErrType: appErr.GetType().String(),
		Message: appErr.GetMessage(),
		Details: appErr.GetDetails(),
	}

	log.Error().
		Str("type", resp.ErrType).
		Str("message", resp.Message).
		Strs("details", resp.Details).
		Msg("app error")

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
