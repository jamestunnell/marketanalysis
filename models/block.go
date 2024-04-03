package models

type Block interface {
	Typed

	Params() Params
	Inputs() Inputs
	Outputs() Outputs

	Warm() bool
	Init() error
	Step()
}

type Input interface {
	Typed
}

type Output interface {
	Typed
}

type Typed interface {
	Type() string
}

type Inputs map[string]Input
type Outputs map[string]Output
