package params

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jamestunnell/marketanalysis/models"
)

type ParamJSON struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func LoadParam(r io.Reader) (models.Param, error) {
	d, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	var paramJSON ParamJSON

	if err = json.Unmarshal(d, &paramJSON); err != nil {
		err = fmt.Errorf("failed to unmarshal param JSON: %w", err)

		return nil, err
	}

	var param models.Param

	switch paramJSON.Type {
	case TypeFloat:
		var f float64

		if err = json.Unmarshal(paramJSON.Value, &f); err != nil {
			return nil, fmt.Errorf("failed to unmarshal as float param: %w", err)
		}

		param = NewFloat(f)
	case TypeInt:
		var i int

		if err = json.Unmarshal(paramJSON.Value, &i); err != nil {
			return nil, fmt.Errorf("failed to unmarshal as int param: %w", err)
		}

		param = NewInt(i)
	}

	if param == nil {
		return nil, fmt.Errorf("unknown param type '%s'", paramJSON.Type)
	}

	return param, nil
}

func StoreParam(p models.Param, w io.Writer) error {
	valueD, err := json.Marshal(p.Value())
	if err != nil {
		return fmt.Errorf("failed to marshal value JSON: %w", err)
	}

	paramJSON := &ParamJSON{
		Type:  p.Type(),
		Value: valueD,
	}

	d, err := json.Marshal(paramJSON)
	if err != nil {
		return fmt.Errorf("failed to marshal param JSON: %w", err)
	}

	_, err = w.Write(d)
	if err != nil {
		return fmt.Errorf("failed to write JSON data: %w", err)
	}

	return nil
}
