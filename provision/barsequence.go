package provision

import "github.com/jamestunnell/marketdata"

type BarSequence interface {
	EachBar(func(bar *marketdata.Bar) error) error
}

type BarSequences interface {
	EachSequence(func(seq BarSequence) error) error
}
