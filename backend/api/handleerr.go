package api

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func handleErr(w http.ResponseWriter, err error, statusCode int) {
	switch statusCode {
	case http.StatusInternalServerError:
		log.Error().Err(err).Int("status", statusCode).Msg("internal failure")
	}

	resp := fmt.Sprintf("{message:\"%v\"}", err)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	_, err = w.Write([]byte(resp))
	if err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
