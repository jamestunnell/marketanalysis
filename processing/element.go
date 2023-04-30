package processing

import (
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
)

type Element interface {
	Type() string
	Params() models.Params
	WarmupPeriod() int
	Initialize() error
	Output() float64
}

type ElementJSON struct {
	Type   string                     `json:"type"`
	Params map[string]json.RawMessage `json:"params"`
}

func MarshalElement(elem Element) ([]byte, error) {
	ps := map[string]json.RawMessage{}
	for name, param := range elem.Params() {
		d, err := param.StoreVal()
		if err != nil {
			return []byte{}, fmt.Errorf("failed to store param '%s': %w", name, err)
		}

		ps[name] = d
	}

	elemJSON := &ElementJSON{
		Type:   elem.Type(),
		Params: ps,
	}

	return json.Marshal(elemJSON)
}

func UnmarshalElement[T Element](d []byte, reg *ElementRegistry[T]) (T, error) {
	var elem T
	var elemJSON ElementJSON

	if err := json.Unmarshal(d, &elemJSON); err != nil {
		err = fmt.Errorf("failed to unmarshal element JSON: %w", err)

		return elem, err
	}

	newElem, found := reg.Get(elemJSON.Type)
	if !found {
		return elem, fmt.Errorf("unknown element type '%s'", elemJSON.Type)
	}

	elem = newElem()

	for name, p := range elem.Params() {
		rawMsg, found := elemJSON.Params[name]
		if !found {
			return elem, &ErrMissingParam{Name: name}
		}

		err := p.LoadVal(rawMsg)
		if err != nil {
			return elem, fmt.Errorf("failed to load value for param '%s': %w", name, err)
		}
	}

	if err := elem.Initialize(); err != nil {
		err = fmt.Errorf("failed to init element: %w", err)

		return elem, err
	}

	return elem, nil
}
