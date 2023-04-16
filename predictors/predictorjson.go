package predictors

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/params"
)

type PredictorJSON struct {
	Type   string                     `json:"type"`
	Params map[string]json.RawMessage `json:"params"`
}

func LoadPredictorFromFile(fpath string) (models.Predictor, error) {
	f, err := os.Open(fpath)
	if err != nil {
		err = fmt.Errorf("failed to open file %s: %w", fpath, err)

		return nil, err
	}

	defer f.Close()

	return LoadPredictor(f)
}

func LoadPredictor(r io.Reader) (models.Predictor, error) {
	d, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	var stratJSON PredictorJSON

	if err = json.Unmarshal(d, &stratJSON); err != nil {
		err = fmt.Errorf("failed to unmarshal Predictor JSON: %w", err)

		return nil, err
	}

	ps := models.Params{}
	for name, rawMsg := range stratJSON.Params {
		p, err := params.LoadParam(bytes.NewReader(rawMsg))
		if err != nil {
			return nil, fmt.Errorf("failed to load param '%s': %w", name, err)
		}

		ps[name] = p
	}

	var newPredictor func(models.Params) (models.Predictor, error)

	switch stratJSON.Type {
	case TypeMACross:
		newPredictor = NewMACross
	case TypePivot:
		newPredictor = NewPivot
	}

	if newPredictor == nil {
		return nil, fmt.Errorf("unknown Predictor type '%s'", stratJSON.Type)
	}

	strat, err := newPredictor(ps)
	if err != nil {
		err = fmt.Errorf("failed to make %s Predictor: %w", stratJSON.Type, err)

		return nil, err
	}

	return strat, nil
}

func StorePredictorToFile(s models.Predictor, fpath string) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fpath, err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	err = StorePredictor(s, w)
	if err != nil {
		return err
	}

	w.Flush()

	return nil
}

func StorePredictor(s models.Predictor, w io.Writer) error {
	ps := map[string]json.RawMessage{}
	for name, param := range s.Params() {
		var buf bytes.Buffer
		if err := params.StoreParam(param, &buf); err != nil {
			return fmt.Errorf("failed to store param '%s': %w", name, err)
		}

		ps[name] = buf.Bytes()
	}

	stratJSON := &PredictorJSON{
		Type:   s.Type(),
		Params: ps,
	}

	d, err := json.Marshal(stratJSON)
	if err != nil {
		return fmt.Errorf("failed to marshal Predictor JSON: %w", err)
	}

	_, err = w.Write(d)
	if err != nil {
		return fmt.Errorf("failed to write JSON data: %w", err)
	}

	return nil
}
