package record

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
	"golang.org/x/exp/maps"
)

type Record struct {
	Recorder    blocks.Recorder
	Inputs      map[string]*blocks.TypedInput[float64]
	InputsAsync map[string]*blocks.TypedInputAsync[float64]
}

const (
	Descr = "Records block outputs"
	Type  = "Record"
)

func (blk *Record) GetType() string {
	return Type
}

func (blk *Record) GetDescription() string {
	return Descr
}

func (blk *Record) GetParams() blocks.Params {
	return blocks.Params{}
}

func (blk *Record) GetInputs() blocks.Inputs {
	ins := blocks.Inputs{}

	for name, in := range blk.Inputs {
		ins[name] = in
	}

	return ins
}

func (blk *Record) GetOutputs() blocks.Outputs {
	return blocks.Outputs{}
}

func (blk *Record) GetWarmupPeriod() int {
	return 0
}

func (blk *Record) IsWarm() bool {
	return true
}

func (blk *Record) Init() error {
	return blk.Recorder.Init(maps.Keys(blk.Inputs))
}

func (blk *Record) Update(cur *models.Bar) {
	vals := map[string]float64{}

	for name, in := range blk.Inputs {
		if in.IsValueSet() {
			vals[name] = in.GetValue()
		}
	}

	// don't record nothing
	if len(vals) == 0 {
		return
	}

	blk.Recorder.Process(cur.Timestamp, vals)
}
