package multitrend

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
)

type MultiTrend struct {
	ins     []*blocks.TypedInput[float64]
	inNames []string
	out     *blocks.TypedOutput[float64]

	numInputs   *blocks.IntParam
	thresh      *blocks.FloatParam
	votesNeeded *blocks.IntParam

	dirs     []models.Direction
	minVotes int
}

const (
	Type  = "MultiTrend"
	Descr = `Aggregates trend inputs (>= 0 for uptrend, < 0 for downtrend)`

	NameVotesNeeded = "votesNeeded"
	NameNumIns      = "numIns"
	NameThresh      = "thresh"
)

func New() blocks.Block {
	return &MultiTrend{
		numInputs:   blocks.NewIntParam(2, blocks.NewGreaterEqual(1)),
		thresh:      blocks.NewFloatParam(0.375, blocks.NewRangeExcl(0.0, 1.0)),
		votesNeeded: blocks.NewIntParam(1, blocks.NewGreaterEqual(1)),
		ins:         []*blocks.TypedInput[float64]{},
		inNames:     []string{},
		out:         blocks.NewTypedOutput[float64](),
		dirs:        []models.Direction{},
		minVotes:    0,
	}
}

func (blk *MultiTrend) GetType() string {
	return Type
}

func (blk *MultiTrend) GetDescription() string {
	return Descr
}

func (blk *MultiTrend) GetParams() blocks.Params {
	return blocks.Params{
		NameNumIns:      blk.numInputs,
		NameVotesNeeded: blk.votesNeeded,
		NameThresh:      blk.thresh,
	}
}

func (blk *MultiTrend) GetInputs() blocks.Inputs {
	ins := blocks.Inputs{}

	for i := 0; i < len(blk.ins); i++ {
		ins[blk.inNames[i]] = blk.ins[i]
	}

	return ins
}

func (blk *MultiTrend) GetOutputs() blocks.Outputs {
	return blocks.Outputs{blocks.NameOut: blk.out}
}

func (blk *MultiTrend) GetWarmupPeriod() int {
	return 0
}

func (blk *MultiTrend) IsWarm() bool {
	return true
}

func (blk *MultiTrend) Init() error {
	if blk.votesNeeded.CurrentVal > blk.numInputs.CurrentVal {
		return commonerrs.NewErrMoreThanMax("votes needed", blk.votesNeeded.CurrentVal, blk.numInputs.CurrentVal)
	}

	numIns := blk.numInputs.CurrentVal
	ins := make([]*blocks.TypedInput[float64], numIns)
	inNames := make([]string, numIns)
	dirs := make([]models.Direction, numIns)

	for i := 0; i < numIns; i++ {
		ins[i] = blocks.NewTypedInput[float64]()
		inNames[i] = fmt.Sprintf("%s%d", blocks.NameIn, i+1)
		dirs[i] = models.DirNone
	}

	blk.dirs = dirs
	blk.ins = ins
	blk.inNames = inNames

	return nil
}

func (blk *MultiTrend) Update(_ *models.Bar, isLast bool) {
	for _, in := range blk.ins {
		if !in.IsValueSet() {
			return
		}
	}

	votesUp := 0
	votesDown := 0

	// update directions for each input
	for i, in := range blk.ins {
		dir := UpdateDirection(in.GetValue(), blk.thresh.CurrentVal, blk.dirs[i])

		switch dir {
		case models.DirUp:
			votesUp++
		case models.DirDown:
			votesDown++
		}

		blk.dirs[i] = dir
	}

	if votesUp >= blk.votesNeeded.CurrentVal {
		blk.out.SetValue(1.0)
	} else if votesDown >= blk.votesNeeded.CurrentVal {
		blk.out.SetValue(-1.0)
	} else {
		blk.out.SetValue(0.0)
	}
}

func UpdateDirection(val, thresh float64, prevDir models.Direction) models.Direction {
	var dir models.Direction

	switch prevDir {
	case models.DirNone:
		if val > thresh {
			dir = models.DirUp
		} else if val < -thresh {
			dir = models.DirDown
		} else {
			dir = models.DirNone
		}
	case models.DirUp:
		if val < -thresh {
			dir = models.DirDown
		} else {
			dir = models.DirUp
		}
	case models.DirDown:
		if val > thresh {
			dir = models.DirUp
		} else {
			dir = models.DirDown
		}
	}

	return dir
}
