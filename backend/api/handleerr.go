package api

import (
	"fmt"
	"net/http"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func handleAppErr(w http.ResponseWriter, appErr *app.Error) {
	var statusCode int

	switch appErr.Code {
	case app.NotFound:
		log.Error().Err(appErr.Err).Msg("app resource not found")

		statusCode = http.StatusNotFound
	case app.ActionFailed:
		log.Error().Err(appErr.Err).Msg("app action failed")

		statusCode = http.StatusInternalServerError
	default:
		log.Error().
			Err(appErr.Err).
			Int("code", appErr.Code).
			Msg("unknown app error")

		statusCode = http.StatusInternalServerError
	}

	resp := fmt.Sprintf("{message:\"%v\"}", appErr.Err)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	_, err := w.Write([]byte(resp))
	if err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
