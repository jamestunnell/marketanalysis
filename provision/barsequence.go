package provision

import "github.com/jamestunnell/marketanalysis/models"

type BarSequence interface {
	EachBar(func(bar *models.Bar) error) error
}

type BarSequences interface {
	EachSequence(func(seq BarSequence) error) error
}
