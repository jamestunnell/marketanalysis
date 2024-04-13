package models

type Model interface {
	GetType() string
	GetName() string
	GetParams() Params
	GetOutputs() Outputs

	Init() error
	Update(*Bar)
}
