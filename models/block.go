package models

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
