package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jamestunnell/marketanalysis/app/backend"
)

func Update[T backend.Resource](
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

	key := mux.Vars(r)[s.GetInfo().KeyName]
	if jsonKey := val.GetKey(); jsonKey != key {
		handleAppErr(w, backend.NewErrInvalidInput("key in request JSON", "does not match key in URL"))

		return
	}

	if appErr := s.Update(r.Context(), val); appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	if postHook != nil {
		postHook(val)
	}

	w.WriteHeader(http.StatusNoContent)
}
