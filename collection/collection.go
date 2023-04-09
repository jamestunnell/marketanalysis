package collection

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rickb777/date/timespan"
	"golang.org/x/exp/slices"
)

type Collection interface {
	Info() *Info
	Timespan() timespan.TimeSpan
	GetBars(ts timespan.TimeSpan) models.Bars
	AddBars(models.Bars) int

	Store(s Store) error
}

type collection struct {
	bars  models.Bars
	info  *Info
	store Store
}

const (
	BarsItemName = "bars.jsonl"
	InfoItemName = "info.json"
)

func Exists(store Store) (bool, error) {
	names, err := store.ItemNames()
	if err != nil {
		err = fmt.Errorf("failed to get store item names: %w", err)

		return false, err
	}

	reqdItems := []string{InfoItemName, BarsItemName}
	for _, reqdName := range reqdItems {
		if !slices.Contains(names, reqdName) {
			return false, nil
		}
	}

	return true, nil
}

func Load(store Store) (Collection, error) {
	d, err := store.LoadItem(InfoItemName)
	if err != nil {
		return nil, fmt.Errorf("failed load info item: %w", err)
	}

	var info Info
	if err = json.Unmarshal(d, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal info: %w", err)
	}

	d, err = store.LoadItem(BarsItemName)
	if err != nil {
		return nil, fmt.Errorf("failed load bars item: %w", err)
	}

	bars, err := models.LoadBars(bytes.NewReader(d))
	if err != nil {
		return nil, fmt.Errorf("failed process bars data: %w", err)
	}

	c := &collection{
		info: &info,
		bars: bars,
	}

	return c, nil
}

func New(info *Info, bars models.Bars) (Collection, error) {
	c := &collection{
		info: info,
		bars: bars,
	}

	return c, nil
}

func (c *collection) Info() *Info {
	return c.info
}

func (c *collection) Timespan() timespan.TimeSpan {
	return c.bars.Timespan()
}

func (c *collection) GetBars(ts timespan.TimeSpan) models.Bars {
	bars := models.Bars{}

	for _, bar := range c.bars {
		if ts.Contains(bar.Timestamp) {
			bars = append(bars, bar)
		}
	}

	return bars
}

func (c *collection) AddBars(bars models.Bars) int {
	added := 0

	for _, bar := range bars {
		found := false
		for _, existingBar := range c.bars {
			if existingBar.Timestamp.Equal(bar.Timestamp) {
				found = true
				break
			}
		}

		if !found {
			c.bars = append(c.bars, bar)

			added++
		}
	}

	return added
}

func (c *collection) Store(s Store) error {
	d, err := json.Marshal(c.info)
	if err != nil {
		return fmt.Errorf("failed to marshal info: %w", err)
	}

	err = s.StoreItem(InfoItemName, d)
	if err != nil {
		return fmt.Errorf("failed store info item: %w", err)
	}

	var b bytes.Buffer

	c.SortBars()

	err = c.bars.Store(&b)
	if err != nil {
		return fmt.Errorf("failed make bars data: %w", err)
	}

	err = s.StoreItem(BarsItemName, b.Bytes())
	if err != nil {
		return fmt.Errorf("failed store bars item: %w", err)
	}

	return nil
}

func (c *collection) SortBars() {
	slices.SortFunc(c.bars, func(a, b *models.Bar) bool {
		return a.Timestamp.Before(b.Timestamp)
	})
}
