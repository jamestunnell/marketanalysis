package collection

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"golang.org/x/exp/slices"
)

type Collection interface {
	Info() *Info
	Dates() []date.Date
	GetBars(date.Date) models.Bars
	AddBars(models.Bars) int
	Store() error
}

type collection struct {
	info     *Info
	index    *DateIndex
	store    Store
	barStore Store
}

const (
	BarsStoreName = "bars"
	InfoItemName  = "info.json"
)

func Exists(store Store) bool {
	if !slices.Contains(store.ItemNames(), InfoItemName) {
		return false
	}

	if !slices.Contains(store.SubstoreNames(), BarsStoreName) {
		return false
	}

	return true
}

func Load(store Store) (Collection, error) {
	infoItem, err := store.Item(InfoItemName)
	if err != nil {
		return nil, fmt.Errorf("failed to get info item: %w", err)
	}

	barStore, err := store.Substore(BarsStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bars substore: %w", err)
	}

	d, err := infoItem.Load()
	if err != nil {
		return nil, fmt.Errorf("failed load info item: %w", err)
	}

	var info Info
	if err = json.Unmarshal(d, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal info: %w", err)
	}

	idx := NewDateIndex(barStore)

	c := &collection{
		info:     &info,
		store:    store,
		barStore: barStore,
		index:    idx,
	}

	return c, nil
}

func New(info *Info, store Store) (Collection, error) {
	if !slices.Contains(store.SubstoreNames(), BarsStoreName) {
		if err := store.MakeSubstore(BarsStoreName); err != nil {
			return nil, fmt.Errorf("failed to make bars substore: %w", err)
		}
	}

	barsStore, err := store.Substore(BarsStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bars substore: %w", err)
	}

	idx := NewDateIndex(barsStore)

	c := &collection{
		info:     info,
		barStore: barsStore,
		store:    store,
		index:    idx,
	}

	return c, nil
}

func (c *collection) Info() *Info {
	return c.info
}

func (c *collection) Dates() []date.Date {
	return c.index.Dates()
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
