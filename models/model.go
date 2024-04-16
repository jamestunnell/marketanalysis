package models

type Model interface {
	GetType() string
	GetName() string
	GetParams() Params
	GetOutputs() Outputs

	Init(Recorder) error
	Update()
}
