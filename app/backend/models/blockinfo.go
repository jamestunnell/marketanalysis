package models

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	m "github.com/jamestunnell/marketanalysis/models"
)

type BlockInfo struct {
	Type       string   `json:"type"`
	Descr      string   `json:"description"`
	Parameters []*Param `json:"parameters"`
	Inputs     []*Input `json:"inputs"`
	Outputs    []*Port  `json:"outputs"`
}

type Param struct {
	Name         string            `json:"name"`
	ValueType    string            `json:"valueType"`
	DefaultValue any               `json:"defaultValue"`
	Constraint   *m.ConstraintInfo `json:"constraint"`
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
		params = append(params, &Param{
			Name:         name,
			ValueType:    p.GetValueType(),
			DefaultValue: p.GetDefaultVal(),
			Constraint:   p.GetConstraintInfo(),
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
		Type:       b.GetType(),
		Descr:      b.GetDescription(),
		Parameters: params,
		Inputs:     ins,
		Outputs:    outs,
	}
}
