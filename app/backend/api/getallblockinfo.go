package api

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/app/backend/models"
	"github.com/jamestunnell/marketanalysis/blocks/registry"
)

type getAllBlockInfo struct{}

func NewGetAllBlockInfo() http.Handler {
	return &getAllBlockInfo{}
}

func (h *getAllBlockInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	infos := []*models.BlockInfo{}

	types := registry.Types()

	slices.Sort(types)

	for _, typ := range types {

		new, _ := registry.Get(typ)

		infos = append(infos, models.NewBlockInfo(new()))
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	p := map[string]any{"blocks": infos}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Warn().Msg("failed to write response")
	}
}
