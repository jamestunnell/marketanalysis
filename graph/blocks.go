package graph

import (
	"errors"
	"fmt"

	gr "github.com/dominikbraun/graph"
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/rs/zerolog/log"
)

type Blocks map[string]blocks.Block

func (blocks Blocks) Init() error {
	for name, block := range blocks {
		if err := block.Init(); err != nil {
			return fmt.Errorf("failed to init block %s: %w", name, err)
		}
	}

	return nil
}

func (blocks Blocks) FindParam(addr *Address) (blocks.Param, bool) {
	block, found := blocks[addr.A]
	if !found {
		return nil, false
	}

	p, found := block.GetParams()[addr.B]
	if !found {
		return nil, false
	}

	return p, true
}

func (blocks Blocks) FindInput(addr *Address) (blocks.Input, bool) {
	block, found := blocks[addr.A]
	if !found {
		return nil, false
	}

	in, found := block.GetInputs()[addr.B]
	if !found {
		return nil, false
	}

	return in, true
}

func (blocks Blocks) FindOutput(addr *Address) (blocks.Output, bool) {
	block, found := blocks[addr.A]
	if !found {
		return nil, false
	}

	out, found := block.GetOutputs()[addr.B]
	if !found {
		return nil, false
	}

	return out, true
}

func (blocks Blocks) Connect(
	conns []*Connection,
) (gr.Graph[string, string], error) {
	g := gr.New(gr.StringHash, gr.Directed(), gr.Acyclic())

	for name, block := range blocks {
		if err := g.AddVertex(name); err != nil {
			return nil, fmt.Errorf("failed to add graph vertex: %w", err)
		}

		for _, out := range block.GetOutputs() {
			out.DisconnectAll()
		}
	}

	for _, c := range conns {
		if err := g.AddEdge(c.Source.A, c.Target.A); err != nil {
			if !errors.Is(err, gr.ErrEdgeAlreadyExists) {
				return nil, fmt.Errorf("failed to add graph edge from %s to %s: %w", c.Source.A, c.Target.A, err)
			}
		}

		output, found := blocks.FindOutput(c.Source)
		if !found {
			return nil, fmt.Errorf("output %s not found", c.Source)
		}

		input, found := blocks.FindInput(c.Target)
		if !found {
			return nil, fmt.Errorf("input %s not found", c.Target)
		}

		output.Connect(input)

		log.Debug().
			Stringer("src", c.Source).
			Stringer("tgt", c.Target).
			Msg("connected pair")
	}

	for blkName, block := range blocks {
		for inName, in := range block.GetInputs() {
			if !in.IsConnected() && !in.IsOptional() {
				err := fmt.Errorf("block %s has unconnected required input '%s'", blkName, inName)

				return nil, err
			}
		}
	}

	return g, nil
}
