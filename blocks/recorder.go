package blocks

import "time"

type Recorder interface {
	Init(valNames []string) error
	Process(t time.Time, vals map[string]float64)
	Finalize() error
}
