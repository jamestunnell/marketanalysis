package predictors

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jamestunnell/marketanalysis/models"
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

	var predJSON PredictorJSON

	if err = json.Unmarshal(d, &predJSON); err != nil {
		err = fmt.Errorf("failed to unmarshal Predictor JSON: %w", err)

		return nil, err
	}

	newPredictor, found := GetNewPredictorFunc(predJSON.Type)
	if !found {
		return nil, fmt.Errorf("unknown Predictor type '%s'", predJSON.Type)
	}

	pred := newPredictor()

	for name, p := range pred.Params() {
		rawMsg, found := predJSON.Params[name]
		if !found {
			return nil, &ErrMissingParam{Name: name}
		}

		err := p.LoadVal(rawMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to load value for param '%s': %w", name, err)
		}
	}

	if err := pred.Initialize(); err != nil {
		err = fmt.Errorf("failed to nit %s Predictor: %w", predJSON.Type, err)

		return nil, err
	}

	return pred, nil
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
		d, err := param.StoreVal()
		if err != nil {
			return fmt.Errorf("failed to store param '%s': %w", name, err)
		}

		ps[name] = d
	}

	predJSON := &PredictorJSON{
		Type:   s.Type(),
		Params: ps,
	}

	d, err := json.Marshal(predJSON)
	if err != nil {
		return fmt.Errorf("failed to marshal Predictor JSON: %w", err)
	}

	_, err = w.Write(d)
	if err != nil {
		return fmt.Errorf("failed to write JSON data: %w", err)
	}

	return nil
}
