package models

type Block interface {
	GetType() string
	GetDescription() string
	GetParams() Params
	GetInputs() Inputs
	GetOutputs() Outputs

	IsWarm() bool

	Init() error
	Update()
}
