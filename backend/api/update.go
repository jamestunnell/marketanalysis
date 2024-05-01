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

	if appErr := s.UpdateFromJSON(r.Context(), key, r.Body); appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
