package models

import (
	"github.com/jamestunnell/marketanalysis/blocks"
)

type BlockInfo struct {
	Type    string   `json:"type"`
	Descr   string   `json:"description"`
	Params  []*Param `json:"params"`
	Inputs  []*Port  `json:"inputs"`
	Outputs []*Port  `json:"outputs"`
}

type Param struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Default any    `json:"default"`
	Limits  []any  `json:"limits"`
}

type Port struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewBlockInfo(b blocks.Block) *BlockInfo {
	params := []*Param{}
	for name, p := range b.GetParams() {
		params = append(params, &Param{
			Name:    name,
			Type:    p.GetType(),
			Default: p.GetDefault(),
			Limits:  p.GetLimits(),
		})
	}

	ins := []*Port{}
	for name, in := range b.GetInputs() {
		ins = append(ins, &Port{Name: name, Type: in.GetType()})
	}

	outs := []*Port{}
	for name, out := range b.GetOutputs() {
		outs = append(outs, &Port{Name: name, Type: out.GetType()})
	}

	return &BlockInfo{
		Type:    b.GetType(),
		Descr:   b.GetDescription(),
		Params:  params,
		Inputs:  ins,
		Outputs: outs,
	}
}
