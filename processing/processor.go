package processing

type Processor interface {
	Element

	WarmUp(vals []float64)
	Update(val float64)
}

var processorRegistry = NewElementRegistry[Processor]()

func MarshalProcessor(p Processor) ([]byte, error) {
	return MarshalElement(p)
}

func UnmarshalProcessor(d []byte) (Processor, error) {
	return UnmarshalElement(d, processorRegistry)
}
