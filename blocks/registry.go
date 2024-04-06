package blocks

import (
	"golang.org/x/exp/maps"
	
	"github.com/jamestunnell/marketanalysis/models"
)

type BlockRegistry struct {
	registered map[string]NewBlockFunc
}

type NewBlockFunc func() models.Block

var registry = NewBlockRegistry()

func Registry() *BlockRegistry {
	return registry
}

func NewBlockRegistry() *BlockRegistry {
	return &BlockRegistry{
		registered: map[string]NewBlockFunc{},
	}
}

func (r *BlockRegistry) Types() []string {
	return maps.Keys(r.registered)
}

func (r *BlockRegistry) Add(typ string, newBlock NewBlockFunc) {
	r.registered[typ] = newBlock
}

func (r *BlockRegistry) Get(typ string) (NewBlockFunc, bool) {
	f, found := r.registered[typ]
	if !found {
		return nil, false
	}

	return f, true
}
