package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/rickb777/date/timespan"
)

type Bars []*Bar

func LoadBarsFromFile(fpath string) (Bars, error) {
	f, err := os.Open(fpath)
	if err != nil {
		err = fmt.Errorf("failed to open file %s: %w", fpath, err)

		return Bars{}, err
	}

	defer f.Close()

	return LoadBars(f)
}

func LoadBars(r io.Reader) (Bars, error) {
	scanner := bufio.NewScanner(r)
	bars := make(Bars, 0)
	line := 1

	for scanner.Scan() {
		var bar Bar

		lineStr := scanner.Text()

		// ignore empty lines
		if lineStr == "" {
			continue
		}

		if err := json.Unmarshal([]byte(lineStr), &bar); err != nil {
			err = fmt.Errorf("failed to unmarshal line %d :%w", line, err)

			return Bars{}, err
		}

		bars = append(bars, &bar)

		line++
	}

	if err := scanner.Err(); err != nil {
		err = fmt.Errorf("scanner failed: %w", err)

		return Bars{}, err
	}

	return bars, nil
}

func (bars Bars) ClosePrices() []float64 {
	closes := make([]float64, len(bars))

	for i, bar := range bars {
		closes[i] = bar.Close
	}

	return closes
}

func (bars Bars) Timespan() timespan.TimeSpan {
	if len(bars) == 0 {
		return timespan.TimeSpan{}
	}

	min := bars[0].Timestamp
	max := bars[0].Timestamp

	for i := 1; i < len(bars); i++ {
		if bars[i].Timestamp.Before(min) {
			min = bars[i].Timestamp
		}

		if bars[i].Timestamp.After(max) {
			max = bars[i].Timestamp
		}
	}

	return timespan.NewTimeSpan(min, max)
}

func (bars Bars) StoreToFile(fpath string) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fpath, err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	return bars.Store(w)
}

func (bars Bars) Store(w io.Writer) error {
	for i, bar := range bars {
		d, err := json.Marshal(bar)
		if err != nil {
			return fmt.Errorf("failed to marshal bar %d: %w", i+1, err)
		}

		d = append(d, byte('\n'))

		_, err = w.Write(d)
		if err != nil {
			return fmt.Errorf("failed to write bar %d: %w", i+1, err)
		}
	}

	return nil
}
