package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jamestunnell/marketanalysis/app/backend"
)

func Delete[T backend.Resource](
	w http.ResponseWriter,
	r *http.Request,
	s backend.Store[T],
	postHook func(string),
) {
	keyVal := mux.Vars(r)[s.GetInfo().KeyName]

	appErr := s.Delete(r.Context(), keyVal)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	if postHook != nil {
		postHook(keyVal)
	}

	w.WriteHeader(http.StatusNoContent)
}
