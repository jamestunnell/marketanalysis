package api

import (
	"encoding/json"
	"net/http"

	"github.com/jamestunnell/marketanalysis/app"
)

func GetAll[T app.Resource](
	w http.ResponseWriter,
	r *http.Request,
	s app.Store[T],
) {
	vals, appErr := s.GetAll(r.Context())
	if appErr != nil {
		handleAppErr(w, appErr)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	p := map[string][]T{s.GetInfo().NamePlural: vals}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		appErr := app.NewErrActionFailed("marshal response JSON", err.Error())

		handleAppErr(w, appErr)

		return
	}
}