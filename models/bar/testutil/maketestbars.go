package testutil

import (
	"fmt"
	"os"

	"github.com/jamestunnell/marketanalysis/models/bar"
)

func MakeTestBars(jsonStr string) ([]*bar.Bar, error) {
	f, err := os.CreateTemp("", "testutil*")

	if err != nil {
		return []*bar.Bar{}, fmt.Errorf("failed to create temp: %w", err)
	}

	fname := f.Name()

	defer os.Remove(fname)

	_, err = f.WriteString(jsonStr)
	if err != nil {
		return []*bar.Bar{}, fmt.Errorf("failed to write JSON string: %w", err)
	}

	if err = f.Close(); err != nil {
		return []*bar.Bar{}, fmt.Errorf("failed to close temp file: %w", err)
	}

	bars, err := bar.LoadBarsFromFile(fname)
	if err != nil {
		return []*bar.Bar{}, fmt.Errorf("failed to load bars from temp file: %w", err)
	}

	return bars, nil
}
