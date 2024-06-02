package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"github.com/rickb777/date/timespan"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type Bars []*Bar

func (bars Bars) Sort() {
	slices.SortFunc(bars, CompareBarsByTimestamp)
}

func (bars Bars) BinarySearch(t time.Time) (int, bool) {
	return slices.BinarySearchFunc(bars, t, CompareBarWithTimestamp)
}

func CompareBarsByTimestamp(a, b *Bar) int {
	return a.Timestamp.Compare(b.Timestamp)
}

func CompareBarWithTimestamp(bar *Bar, t time.Time) int {
	return bar.Timestamp.Compare(t)
}

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

func (bars Bars) Last() *Bar {
	return sliceutils.Last(bars)
}

func (bars Bars) LastN(n int) Bars {
	return sliceutils.LastN(bars, n)
}

func (bars Bars) NextN(index, n int) Bars {
	if index < 0 || index >= len(bars) {
		return Bars{}
	}

	a := index + 1
	b := a + n
	if b > len(bars) {
		b = len(bars)
	}

	return bars[a:b]
}

func (bars Bars) Localize() {
	for _, b := range bars {
		b.Localize()
	}
}

func (bars Bars) Timestamps() []time.Time {
	return sliceutils.Map(bars, func(b *Bar) time.Time {
		return b.Timestamp
	})
}

func (bars Bars) ClosePrices() []float64 {
	return sliceutils.Map(bars, func(b *Bar) float64 {
		return b.Close
	})
}

func (bars Bars) TimeSpan() timespan.TimeSpan {
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

	if err = bars.Store(w); err != nil {
		return err
	}

	w.Flush()

	return nil
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
