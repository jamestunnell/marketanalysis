package processing

import (
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/processing/sources"
)

type Source interface {
	Element

	WarmUp(bars models.Bars) error
	Update(bar *models.Bar)
}

var sourceRegistry = NewElementRegistry[Source]()

func init() {
	sourceRegistry.Add(
		sources.TypeCandlestick,
		func() Source { return sources.NewCandlestick() },
	)

	sourceRegistry.Add(
		sources.TypeHeikinAshi,
		func() Source { return sources.NewHeikinAshi() },
	)
}

func MarshalSource(s Source) ([]byte, error) {
	return MarshalElement(s)
}

func UnmarshalSource(d []byte) (Source, error) {
	return UnmarshalElement(d, sourceRegistry)
}
