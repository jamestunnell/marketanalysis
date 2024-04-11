package models

import (
	"fmt"

	"github.com/dominikbraun/graph"
)

type Block interface {
	GetType() string
	GetDescription() string
	GetParams() Params
	GetInputs() Inputs
	GetOutputs() Outputs

	IsWarm() bool

	Init() error
	Update(*Bar)
}

type Blocks map[string]Block

type NewBlockFunc func() Block

type BlockRegistry interface {
	Types() []string
	Add(typ string, newBlock NewBlockFunc)
	Get(typ string) (NewBlockFunc, bool)
}

func (blocks Blocks) Init() error {
	for name, block := range blocks {
		if err := block.GetParams().Validate(); err != nil {
			return fmt.Errorf("block %s has invalid param(s): %w", name, err)
		}

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

func (blocks Blocks) Connect(conns Connections) ([]string, error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())

	for name, block := range blocks {
		if err := g.AddVertex(name); err != nil {
			return []string{}, fmt.Errorf("failed to add graph vertex: %w", err)
		}

		for _, out := range block.GetOutputs() {
			out.DisconnectAll()
		}
	}

	err := conns.EachPair(func(src, tgt *Address) error {
		output, found := blocks.FindOutput(src)
		if !found {
			return fmt.Errorf("output %s not found", src)
		}

		input, found := blocks.FindInput(tgt)
		if !found {
			return fmt.Errorf("input %s not found", tgt)
		}

		output.Connect(input)

		if err := g.AddEdge(src.A, tgt.A); err != nil {
			return fmt.Errorf("failed to add graph edge: %w", err)
		}

		return nil
	})
	if err != nil {
		return []string{}, err
	}

	order, err := graph.TopologicalSort(g)
	if err != nil {
		return []string{}, fmt.Errorf("topological sort failed: %w", err)
	}

	return order, nil
}
