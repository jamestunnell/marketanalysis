package collection

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"golang.org/x/exp/slices"
)

type Collection interface {
	Info() *Info
	TimeSpan() timespan.TimeSpan
	LoadBars(timespan.TimeSpan) (models.Bars, error)
	StoreBars(models.Bars) error
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
	d, err := store.LoadItem(InfoItemName)
	if err != nil {
		return nil, fmt.Errorf("failed to load info item: %w", err)
	}

	var info Info
	if err = json.Unmarshal(d, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal info: %w", err)
	}

	barStore, err := store.Substore(BarsStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bars substore: %w", err)
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
	if !slices.Contains(store.ItemNames(), InfoItemName) {
		d, err := json.Marshal(info)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal info: %w", err)
		}

		if err := store.StoreItem(InfoItemName, d); err != nil {
			return nil, fmt.Errorf("failed to store info item: %w", err)
		}
	}

	var barStore Store

	var err error

	if !slices.Contains(store.SubstoreNames(), BarsStoreName) {
		barStore, err = store.MakeSubstore(BarsStoreName)
		if err != nil {
			return nil, fmt.Errorf("failed to make bars substore: %w", err)
		}
	} else {
		barStore, err = store.Substore(BarsStoreName)
		if err != nil {
			return nil, fmt.Errorf("failed to get bars substore: %w", err)
		}
	}

	idx := NewDateIndex(barStore)

	c := &collection{
		info:     info,
		barStore: barStore,
		store:    store,
		index:    idx,
	}

	return c, nil
}

func (c *collection) Info() *Info {
	return c.info
}

func (c *collection) TimeSpan() timespan.TimeSpan {
	return c.index.TimeSpan()
}

func (c *collection) loadBarsForDate(d date.Date) (models.Bars, error) {
	itemName, found := c.index.FindItem(d)
	if !found {
		return models.Bars{}, nil
	}

	barsData, err := c.barStore.LoadItem(itemName)
	if err != nil {
		err = fmt.Errorf("failed to load bars item '%s': %w", itemName, err)

		return models.Bars{}, err
	}

	bars, err := models.LoadBars(bytes.NewReader(barsData))
	if err != nil {
		return nil, fmt.Errorf("failed to load bars from data: %w", err)
	}

	return bars, nil
}

func (c *collection) LoadBars(ts timespan.TimeSpan) (models.Bars, error) {
	bars := models.Bars{}

	for cur := ts.Start(); cur.Before(ts.End()); cur = cur.AddDate(0, 0, 1) {
		d := date.NewAt(cur)

		dayBars, err := c.loadBarsForDate(d)
		if err != nil {
			dStr := d.Format(date.RFC3339)
			err = fmt.Errorf("failed to laod bars on date %s: %w", dStr, err)

			return models.Bars{}, err
		}

		if len(dayBars) == 0 {
			continue
		}

		if !ts.Contains(dayBars[0].Timestamp) || !ts.Contains(dayBars.Last().Timestamp) {
		} else {
			dayBars = sliceutils.Where(dayBars, func(b *models.Bar) bool {
				return ts.Contains(b.Timestamp)
			})
		}

		bars = append(bars, dayBars...)
	}

	return bars, nil
}

func BarsItemName(sym string, d date.Date) string {
	const fmtStr = "%s_%s.jsonl"

	return fmt.Sprintf(fmtStr, sym, d.Format(date.RFC3339))
}

func (c *collection) StoreBars(bars models.Bars) error {
	// separate by date
	byDate := map[date.Date]models.Bars{}
	for _, b := range bars {
		d := b.Date()

		if dateBars, found := byDate[d]; found {
			byDate[d] = append(dateBars, b)
		} else {
			byDate[d] = models.Bars{b}
		}
	}

	//store each set of bars in an item
	for d, bars := range byDate {
		slices.SortFunc(bars, func(a, b *models.Bar) bool {
			return a.Timestamp.Before(b.Timestamp)
		})

		var buf bytes.Buffer

		if err := bars.Store(&buf); err != nil {
			return fmt.Errorf("failed to make bar data: %w", err)
		}

		name := BarsItemName(c.info.Symbol, d)
		if err := c.barStore.StoreItem(name, buf.Bytes()); err != nil {
			return fmt.Errorf("failed to store bar date: %w", err)
		}
	}

	return nil
}
