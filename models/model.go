package models

type Model interface {
	GetType() string
	GetName() string
	GetParams() Params
	GetOutputs() Outputs
	GetWarmupPeriod() int

	Init(Recorder) error
	Update(*Bar)
}
