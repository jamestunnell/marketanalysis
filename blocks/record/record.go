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

	for name, in := range blk.InputsAsync {
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
	inputNames := append(maps.Keys(blk.Inputs), maps.Keys(blk.InputsAsync)...)

	return blk.Recorder.Init(inputNames)
}

func (blk *Record) Update(cur *models.Bar, isLast bool) {
	tvs := map[string]models.TimeValue[float64]{}

	for name, in := range blk.Inputs {
		if in.IsValueSet() {
			tvs[name] = models.NewTimeValue(cur.Timestamp, in.GetValue())
		}
	}

	for name, in := range blk.InputsAsync {
		if in.IsValueSet() {
			tvs[name] = models.NewTimeValue(in.GetTime(), in.GetValue())
		}
	}

	// don't record nothing
	if len(tvs) == 0 {
		return
	}

	blk.Recorder.Process(tvs)
}
