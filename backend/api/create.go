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
	val, err := LoadRequestJSON[T](r)
	if err != nil {
		handleAppErr(w, app.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	if appErr := s.Create(r.Context(), val); appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
