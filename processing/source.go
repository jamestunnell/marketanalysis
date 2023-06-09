package processing

import (
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/processing/sources"
)

type Source interface {
	Element

	Update(bar *models.Bar)
}

var sourceRegistry = NewElementRegistry[Source]()

func init() {
	sourceRegistry.Add(
		sources.TypeCandlestick,
		func() Source { return sources.NewCandlestick() },
	)

	sourceRegistry.Add(
		sources.TypeDMI,
		func() Source { return sources.NewDMI() },
	)

	sourceRegistry.Add(
		sources.TypeEMV,
		func() Source { return sources.NewEMV() },
	)

	sourceRegistry.Add(
		sources.TypeHeikinAshi,
		func() Source { return sources.NewHeikinAshi() },
	)

	sourceRegistry.Add(
		sources.TypeTrueRange,
		func() Source { return sources.NewTrueRange() },
	)

}

func MarshalSource(s Source) ([]byte, error) {
	return MarshalElement(s)
}

func UnmarshalSource(d []byte) (Source, error) {
	return UnmarshalElement(d, sourceRegistry)
}
