package api

import (
	"encoding/json"
	"net/http"

	"github.com/jamestunnell/marketanalysis/app"
)

func Create[T app.Resource](
	w http.ResponseWriter,
	r *http.Request,
	s app.Store[T],
) {
	val := app.NewResource[T]()
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		handleAppErr(w, app.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	if appErr := s.Create(r.Context(), val); appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
