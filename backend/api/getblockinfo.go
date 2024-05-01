package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app"
	"github.com/jamestunnell/marketanalysis/blocks/registry"
)

type getBlockInfo struct{}

func NewGetBlockInfo() http.Handler {
	return &getBlockInfo{}
}

func (h *getBlockInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	typ := mux.Vars(r)["type"]

	new, found := registry.Get(typ)
	if !found {
		appErr := app.NewNotFoundError(fmt.Sprintf("block with type '%s'", typ))

		handleAppErr(w, appErr)

		return
	}

	info := NewBlockInfo(new())

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Warn().Msg("failed to write response")
	}
}
