package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app"
)

func Update[T any](
	w http.ResponseWriter,
	r *http.Request,
	s app.Store[T],
) {
	key := mux.Vars(r)[s.RDef().KeyName]

	val, err := LoadRequestJSON[T](r)
	if err != nil {
		handleAppErr(w, app.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	if jsonKey := s.RDef().GetKey(val); jsonKey != key {
		handleAppErr(w, app.NewErrInvalidInput("key in request JSON", "does not match key in URL"))

		return
	}

	if appErr := s.Update(r.Context(), val); appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
