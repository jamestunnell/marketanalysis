package models

import (
	"github.com/jamestunnell/marketanalysis/blocks"
)

type BlockInfo struct {
	Type        string            `json:"type"`
	Descr       string            `json:"description"`
	Params      map[string]*Param `json:"params"`
	InputTypes  map[string]string `json:"inputTypes"`
	OutputTypes map[string]string `json:"outputTypes"`
}

type Param struct {
	Default any            `json:"default"`
	Schema  map[string]any `json:"schema"`
}

func NewBlockInfo(b blocks.Block) *BlockInfo {
	params := map[string]*Param{}

	for name, p := range b.GetParams() {
		params[name] = &Param{
			Default: p.GetDefault(),
			Schema:  p.GetSchema(),
		}
	}

	inTypes := map[string]string{}
	for name, in := range b.GetInputs() {
		inTypes[name] = in.GetType()
	}

	outTypes := map[string]string{}
	for name, out := range b.GetOutputs() {
		outTypes[name] = out.GetType()
	}

	return &BlockInfo{
		Type:        b.GetType(),
		Descr:       b.GetDescription(),
		Params:      params,
		InputTypes:  inTypes,
		OutputTypes: outTypes,
	}
}
