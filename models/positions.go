package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Positions []*Position

type PositionsAnalysis struct {
	Winning, TotalPL float64
}

func (ps Positions) Localize() {
	for _, ps := range ps {
		ps.Localize()
	}
}

func (ps Positions) Analyze() *PositionsAnalysis {
	totalPL := 0.0
	nWinning := 0

	for _, pos := range ps {
		if pos.ClosedPL > 0.0 {
			nWinning++
		}

		totalPL += pos.ClosedPL
	}

	winning := float64(nWinning) / float64(len(ps))

	return &PositionsAnalysis{
		Winning: winning,
		TotalPL: totalPL,
	}
}

func (ps Positions) StoreToFile(fpath string) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fpath, err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	if err = ps.Store(w); err != nil {
		return err
	}

	w.Flush()

	return nil
}

func (ps Positions) Store(w io.Writer) error {
	for i, p := range ps {
		d, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("failed to marshal position %d: %w", i+1, err)
		}

		d = append(d, byte('\n'))

		_, err = w.Write(d)
		if err != nil {
			return fmt.Errorf("failed to write position %d: %w", i+1, err)
		}
	}

	return nil
}
