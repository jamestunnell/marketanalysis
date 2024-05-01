package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func LoadRequestJSON[T any](r *http.Request) (*T, error) {
	var val T

	if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	return &val, nil
}
