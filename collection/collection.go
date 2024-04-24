package collection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type collection struct {
	info     *models.CollectionInfo
	index    *DateIndex
	store    Store
	barStore Store
	loc      *time.Location
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

func LoadFromDir(dir string) (models.Collection, error) {
	store, err := NewDirStore(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to make dir store: %w", err)
	}

	return Load(store)
}

func Load(store Store) (models.Collection, error) {
	d, err := store.LoadItem(InfoItemName)
	if err != nil {
		return nil, fmt.Errorf("failed to load info item: %w", err)
	}

	var info models.CollectionInfo
	if err = json.Unmarshal(d, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal info: %w", err)
	}

	loc, err := time.LoadLocation(info.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location from time zone '%s': %w", info.TimeZone, err)
	}

	barStore, err := store.Substore(BarsStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bars substore: %w", err)
	}

	idx := NewDateIndex(barStore)

	log.Debug().Interface("info", info).Msg("loaded collection")

	c := &collection{
		info:     &info,
		store:    store,
		barStore: barStore,
		index:    idx,
		loc:      loc,
	}

	return c, nil
}

func New(info *models.CollectionInfo, store Store) (models.Collection, error) {
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

	loc, err := time.LoadLocation(info.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location from time zone '%s': %w", info.TimeZone, err)
	}

	idx := NewDateIndex(barStore)

	c := &collection{
		info:     info,
		barStore: barStore,
		store:    store,
		index:    idx,
		loc:      loc,
	}

	return c, nil
}

func (c *collection) GetInfo() *models.CollectionInfo {
	return c.info
}

func (c *collection) GetLocation() *time.Location {
	return c.loc
}

func (c *collection) GetFirstDate() date.Date {
	return c.index.first
}

func (c *collection) GetLastDate() date.Date {
	return c.index.last
}

func (c *collection) IsEmpty() bool {
	return c.index.Empty()
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

func (c *collection) LoadBars(start, endIncl date.Date) (models.Bars, error) {
	bars := models.Bars{}

	for cur := start; !cur.After(endIncl); cur = cur.Add(1) {
		dayBars, err := c.loadBarsForDate(cur)
		if err != nil {
			err = fmt.Errorf("failed to load bars on date %s: %w", cur, err)

			return models.Bars{}, err
		}

		if len(dayBars) == 0 {
			continue
		}

		// if !ts.Contains(dayBars[0].Timestamp) || !ts.Contains(dayBars.Last().Timestamp) {
		// } else {
		// 	dayBars = sliceutils.Where(dayBars, func(b *models.Bar) bool {
		// 		return ts.Contains(b.Timestamp)
		// 	})
		// }

		bars = append(bars, dayBars...)
	}

	log.Debug().
		Stringer("start", start).
		Stringer("endIncl", endIncl).
		Msgf("loaded %d bars", len(bars))

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

		c.index.AddItem(name, d)
	}

	c.index.Update()

	return nil
}
