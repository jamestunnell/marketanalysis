package blocks

import (
	"github.com/jamestunnell/marketanalysis/models"
	"golang.org/x/exp/maps"
)

type Record struct {
	Recorder models.Recorder
	Inputs   map[string]*models.TypedInput[float64]
}

const (
	DescrRecord = "Records block outputs"
	TypeRecord  = "Record"
)

func (blk *Record) GetType() string {
	return TypeRecord
}

func (blk *Record) GetDescription() string {
	return DescrRecord
}

func (blk *Record) GetParams() models.Params {
	return models.Params{}
}

func (blk *Record) GetInputs() models.Inputs {
	ins := models.Inputs{}

	for name, in := range blk.Inputs {
		ins[name] = in
	}

	return ins
}

func (blk *Record) GetOutputs() models.Outputs {
	return models.Outputs{}
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
		if in.IsSet() {
			vals[name] = in.Get()
		}
	}

	// don't record nothing
	if len(vals) == 0 {
		return
	}

	blk.Recorder.Record(cur.Timestamp, vals)
}
