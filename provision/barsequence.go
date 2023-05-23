package provision

import "github.com/jamestunnell/marketanalysis/models"

type BarSequence interface {
	Initialize() error
	EachBar(func(bar *models.Bar))
}

type BarSequences interface {
	Initialize() error
	EachSequence(func(seq BarSequence))
}
