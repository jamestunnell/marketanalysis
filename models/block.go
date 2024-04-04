package models

type Block interface {
	Type() string

	Params() Params
	Inputs() Inputs
	Outputs() Outputs

	Warm() bool
	Init() error
	Step()
}
