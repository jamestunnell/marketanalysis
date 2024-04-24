package blocks

import (
	"golang.org/x/exp/maps"

	"github.com/jamestunnell/marketanalysis/models"
)

type registry struct {
	blocks map[string]models.NewBlockFunc
}

type NewBlockFunc func() models.Block

var reg = &registry{
	blocks: map[string]models.NewBlockFunc{},
}

func Registry() models.BlockRegistry {
	return reg
}

func (r *registry) Types() []string {
	return maps.Keys(r.blocks)
}

func (r *registry) Add(typ string, newBlock models.NewBlockFunc) {
	r.blocks[typ] = newBlock
}

func (r *registry) Get(typ string) (models.NewBlockFunc, bool) {
	f, found := r.blocks[typ]
	if !found {
		return nil, false
	}

	return f, true
}
