package api

import "net/http"

type status struct {
}

func NewStatus() http.Handler {
	return &status{}
}

func (h *status) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
