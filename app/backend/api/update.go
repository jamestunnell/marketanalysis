package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app"
)

func Update[T app.Resource](
	w http.ResponseWriter,
	r *http.Request,
	s app.Store[T],
) {
	val := app.NewResource[T]()
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		handleAppErr(w, app.NewErrInvalidInput("request JSON", err.Error()))

		return
	}

	key := mux.Vars(r)[s.GetInfo().KeyName]
	if jsonKey := val.GetKey(); jsonKey != key {
		handleAppErr(w, app.NewErrInvalidInput("key in request JSON", "does not match key in URL"))

		return
	}

	if appErr := s.Update(r.Context(), val); appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
