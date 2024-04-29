package blocks

import "github.com/jamestunnell/marketanalysis/models"

type Block interface {
	GetType() string
	GetDescription() string
	GetParams() Params
	GetInputs() Inputs
	GetOutputs() Outputs
	GetWarmupPeriod() int

	IsWarm() bool

	Init() error
	Update(*models.Bar)
}
