package blocks

import (
	"github.com/jamestunnell/marketanalysis/models"
)

type Recorder interface {
	Init(valNames []string) error
	Process(map[string]models.TimeValue[float64])
	Finalize() error
}
