package processing

import (
	"github.com/jamestunnell/marketanalysis/processing/processors"
)

type Processor interface {
	Element

	WarmUp(vals []float64)
	Update(val float64)
}

var processorRegistry = NewElementRegistry[Processor]()

func init() {
	processorRegistry.Add(
		processors.TypeDiff,
		func() Processor { return processors.NewDiff() },
	)

	processorRegistry.Add(
		processors.TypeEMA,
		func() Processor { return processors.NewEMA() },
	)

	processorRegistry.Add(
		processors.TypeMADiff,
		func() Processor { return processors.NewMADiff() },
	)

	processorRegistry.Add(
		processors.TypeMAOrder,
		func() Processor { return processors.NewMAOrder() },
	)

	processorRegistry.Add(
		processors.TypeSMA,
		func() Processor { return processors.NewSMA() },
	)
}

func MarshalProcessor(p Processor) ([]byte, error) {
	return MarshalElement(p)
}

func UnmarshalProcessor(d []byte) (Processor, error) {
	return UnmarshalElement(d, processorRegistry)
}
