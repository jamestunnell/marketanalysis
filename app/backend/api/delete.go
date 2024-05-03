package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jamestunnell/marketanalysis/app"
)

func Delete[T app.Resource](
	w http.ResponseWriter,
	r *http.Request,
	s app.Store[T],
) {
	keyVal := mux.Vars(r)[s.GetInfo().KeyName]

	appErr := s.Delete(r.Context(), keyVal)
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
