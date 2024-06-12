package api

import (
	"encoding/json"
	"net/http"

	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func handleAppErr(w http.ResponseWriter, appErr backend.Error) {
	var statusCode int
	var title string

	switch appErr.GetType() {
	case backend.NotFound:
		statusCode = http.StatusNotFound
		title = "Not Found"
	case backend.InvalidInput:
		statusCode = http.StatusBadRequest
		title = "Invalid Input"
	case backend.ActionFailed:
		statusCode = http.StatusInternalServerError
		title = "Action Failed"
	default:
		log.Error().Msgf("app error type %s is unknown", appErr.GetType())

		statusCode = http.StatusInternalServerError
		title = "Unknown"
	}

	resp := &ErrorResponse{
		Title:   title,
		Message: appErr.GetMessage(),
		Details: appErr.GetDetails(),
	}

	log.Error().
		Str("errType", appErr.GetType().String()).
		Str("message", resp.Message).
		Strs("details", resp.Details).
		Msg("app error")

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
