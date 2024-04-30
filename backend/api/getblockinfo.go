package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

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
		err := fmt.Errorf("block with type '%s' not found", typ)

		handleErr(w, err, http.StatusNotFound)

		return
	}

	info := NewBlockInfo(new())

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Warn().Msg("failed to write response")
	}
}
