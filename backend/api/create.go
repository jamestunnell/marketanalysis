package api

import (
	"net/http"

	"github.com/jamestunnell/marketanalysis/app"
)

func Create[T any](
	w http.ResponseWriter,
	r *http.Request,
	s app.Store[T],
) {
	if appErr := s.CreateFromJSON(r.Context(), r.Body); appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
