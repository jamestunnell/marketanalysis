package blocks

import (
	"errors"
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/rs/zerolog/log"
)

type Blocks map[string]Block

func (blocks Blocks) Init() error {
	for name, block := range blocks {
		if err := block.Init(); err != nil {
			return fmt.Errorf("failed to init block %s: %w", name, err)
		}
	}

	return nil
}

func (blocks Blocks) FindParam(addr *Address) (Param, bool) {
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

func (blocks Blocks) FindInput(addr *Address) (Input, bool) {
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

func (blocks Blocks) FindOutput(addr *Address) (Output, bool) {
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

func (blocks Blocks) Connect(conns Connections) (graph.Graph[string, string], error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())

	for name, block := range blocks {
		if err := g.AddVertex(name); err != nil {
			return nil, fmt.Errorf("failed to add graph vertex: %w", err)
		}

		for _, out := range block.GetOutputs() {
			out.DisconnectAll()
		}
	}

	err := conns.EachPair(func(src, dest *Address) error {
		if err := g.AddEdge(src.A, dest.A); err != nil {
			if !errors.Is(err, graph.ErrEdgeAlreadyExists) {
				return fmt.Errorf("failed to add graph edge: %w", err)
			}
		}

		output, found := blocks.FindOutput(src)
		if !found {
			return fmt.Errorf("output %s not found", src)
		}

		input, found := blocks.FindInput(dest)
		if !found {
			return fmt.Errorf("input %s not found", dest)
		}

		output.Connect(input)

		log.Debug().
			Stringer("out", src).
			Stringer("in", dest).
			Msg("connected pair")

		return nil
	})
	if err != nil {
		return nil, err
	}

	return g, nil
}
