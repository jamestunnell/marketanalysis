package strategies

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

type StrategyJSON struct {
	Type   string                     `json:"type"`
	Params map[string]json.RawMessage `json:"params"`
}

func LoadStrategyFromFile(fpath string) (models.Strategy, error) {
	f, err := os.Open(fpath)
	if err != nil {
		err = fmt.Errorf("failed to open file %s: %w", fpath, err)

		return nil, err
	}

	defer f.Close()

	return LoadStrategy(f)
}

func LoadStrategy(r io.Reader) (models.Strategy, error) {
	d, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	var stratJSON StrategyJSON

	if err = json.Unmarshal(d, &stratJSON); err != nil {
		err = fmt.Errorf("failed to unmarshal strategy JSON: %w", err)

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

	var newStrategy func(models.Params) (models.Strategy, error)

	switch stratJSON.Type {
	case TypeTrendFollower:
		newStrategy = NewTrendFollower
	case TypeScalper:
		newStrategy = NewScalper
	}

	if newStrategy == nil {
		return nil, fmt.Errorf("unknown strategy type '%s'", stratJSON.Type)
	}

	strat, err := newStrategy(ps)
	if err != nil {
		err = fmt.Errorf("failed to make %s strategy: %w", stratJSON.Type, err)

		return nil, err
	}

	return strat, nil
}

func StoreStrategyToFile(s models.Strategy, fpath string) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fpath, err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	err = StoreStrategy(s, w)
	if err != nil {
		return err
	}

	w.Flush()

	return nil
}

func StoreStrategy(s models.Strategy, w io.Writer) error {
	ps := map[string]json.RawMessage{}
	for name, param := range s.Params() {
		var buf bytes.Buffer
		if err := params.StoreParam(param, &buf); err != nil {
			return fmt.Errorf("failed to store param '%s': %w", name, err)
		}

		ps[name] = buf.Bytes()
	}

	stratJSON := &StrategyJSON{
		Type:   s.Type(),
		Params: ps,
	}

	d, err := json.Marshal(stratJSON)
	if err != nil {
		return fmt.Errorf("failed to marshal strategy JSON: %w", err)
	}

	_, err = w.Write(d)
	if err != nil {
		return fmt.Errorf("failed to write JSON data: %w", err)
	}

	return nil
}
