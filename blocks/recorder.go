package blocks

import "time"

type Recorder interface {
	Init(valNames []string) error
	Record(t time.Time, vals map[string]float64)
	Flush()
}
