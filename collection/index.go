package collection

import (
	"regexp"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rickb777/date"
	"golang.org/x/exp/slices"
)

type DateIndex struct {
	store   Store
	entries []*DateIndexEntry
}

type DateIndexEntry struct {
	Date     date.Date
	ItemName string
}

var matchDate = regexp.MustCompile(`[\d]{4}-[\d]{2}-[\d]{2}`)

func NewDateIndex(s Store) *DateIndex {
	idx := &DateIndex{
		store:   s,
		entries: []*DateIndexEntry{},
	}

	idx.Update()

	return idx
}

func (idx *DateIndex) Dates() []date.Date {
	return sliceutils.Map(idx.entries, func(e *DateIndexEntry) date.Date {
		return e.Date
	})
}

func (idx *DateIndex) Update() {
	names := idx.store.ItemNames()
	entries := []*DateIndexEntry{}

	for _, name := range names {
		d, ok := ExtractDate(name)
		if ok {
			e := &DateIndexEntry{
				Date:     d,
				ItemName: name,
			}
			entries = append(entries, e)
		}
	}

	slices.SortFunc(entries, func(a, b *DateIndexEntry) bool {
		return a.Date.Before(b.Date)
	})
}

func ExtractDate(s string) (date.Date, bool) {
	dateStr := matchDate.FindString(s)
	if dateStr == "" {
		return date.Date{}, false
	}

	d, err := date.ParseISO(dateStr)
	if err != nil {
		return date.Date{}, false
	}

	return d, true
}
