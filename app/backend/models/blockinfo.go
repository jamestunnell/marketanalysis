package models

import (
	"github.com/jamestunnell/marketanalysis/blocks"
)

type BlockInfo struct {
	Type    string   `json:"type"`
	Descr   string   `json:"description"`
	Params  []*Param `json:"params"`
	Inputs  []*Input `json:"inputs"`
	Outputs []*Port  `json:"outputs"`
}

type Param struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Default    any         `json:"default"`
	Constraint *Constraint `json:"constraint"`
}

type Constraint struct {
	Type   string `json:"type"`
	Limits []any  `json:"limits"`
}

type Port struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Async bool   `json:"async"`
}

type Input struct {
	*Port
	Optional bool `json:"optional"`
}

func NewBlockInfo(b blocks.Block) *BlockInfo {
	params := []*Param{}
	for name, p := range b.GetParams() {
		constraint := &Constraint{
			Type:   p.GetConstraint().GetType().String(),
			Limits: p.GetConstraint().GetLimits(),
		}

		params = append(params, &Param{
			Name:       name,
			Type:       p.GetValueType(),
			Default:    p.GetDefaultVal(),
			Constraint: constraint,
		})
	}

	ins := []*Input{}
	for name, in := range b.GetInputs() {
		port := &Port{
			Name:  name,
			Type:  in.GetType(),
			Async: in.IsAsynchronous(),
		}

		ins = append(ins, &Input{
			Port:     port,
			Optional: in.IsOptional(),
		})
	}

	outs := []*Port{}
	for name, out := range b.GetOutputs() {
		port := &Port{
			Name: name,
			Type: out.GetType(),
		}

		outs = append(outs, port)
	}

	return &BlockInfo{
		Type:    b.GetType(),
		Descr:   b.GetDescription(),
		Params:  params,
		Inputs:  ins,
		Outputs: outs,
	}
}
