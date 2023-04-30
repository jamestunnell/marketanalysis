package processing

import "github.com/jamestunnell/marketanalysis/models"

type Source interface {
	Element

	WarmUp(bars models.Bars) error
	Update(bar *models.Bar)
}

var sourceRegistry = NewElementRegistry[Source]()

func MarshalSource(s Source) ([]byte, error) {
	return MarshalElement(s)
}

func UnmarshalSource(d []byte) (Source, error) {
	return UnmarshalElement(d, sourceRegistry)
}
