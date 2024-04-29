package models

import (
	"github.com/jamestunnell/marketanalysis/blocks"
)

type BlockInfo struct {
	Type    string      `json:"type"`
	Descr   string      `json:"description"`
	Params  []*ParamDef `json:"params"`
	Inputs  []*NameType `json:"inputs"`
	Outputs []*NameType `json:"outputs"`
}

type ParamDef struct {
	Name    string         `json:"name"`
	Default any            `json:"default"`
	Schema  map[string]any `json:"schema"`
}

func MakeBlockInfo(b blocks.Block) *BlockInfo {
	params := b.GetParams()
	paramDefs := []*ParamDef{}

	for name, p := range params {
		pd := &ParamDef{
			Name:    name,
			Default: p.GetDefault(),
			Schema:  p.GetSchema(),
		}

		paramDefs = append(paramDefs, pd)
	}

	inNameTypes := []*NameType{}
	for name, in := range b.GetInputs() {
		nt := &NameType{Name: name, Type: in.GetType()}

		inNameTypes = append(inNameTypes, nt)
	}

	outNameTypes := []*NameType{}
	for name, out := range b.GetOutputs() {
		nt := &NameType{Name: name, Type: out.GetType()}

		outNameTypes = append(outNameTypes, nt)
	}

	return &BlockInfo{
		Type:    b.GetType(),
		Descr:   b.GetDescription(),
		Params:  paramDefs,
		Inputs:  inNameTypes,
		Outputs: outNameTypes,
	}
}
