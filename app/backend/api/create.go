package api

import (
	"encoding/json"
	"net/http"

	"github.com/jamestunnell/marketanalysis/app/backend"
)

func Create[T backend.Resource](
	w http.ResponseWriter,
	r *http.Request,
	s backend.Store[T],
	postHook func(T),
) {
	val := backend.NewResource[T]()
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		handleAppErr(w, backend.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	if appErr := s.Create(r.Context(), val); appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	if postHook != nil {
		postHook(val)
	}

	w.WriteHeader(http.StatusNoContent)
}
