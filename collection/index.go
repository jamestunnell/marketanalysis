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
	first   date.Date
	last    date.Date
}

type DateIndexEntry struct {
	Date     date.Date
	ItemName string
}

var (
	dummyTgt  = &DateIndexEntry{Date: date.Today()}
	matchDate = regexp.MustCompile(`[\d]{4}-[\d]{2}-[\d]{2}`)
)

func NewDateIndex(s Store) *DateIndex {
	idx := &DateIndex{
		store:   s,
		entries: []*DateIndexEntry{},
		first:   date.Date{},
		last:    date.Date{},
	}

	idx.Update()

	return idx
}

func (idx *DateIndex) Empty() bool {
	return len(idx.entries) == 0
}

func (idx *DateIndex) LastDate() date.Date {
	return idx.last
}

func (idx *DateIndex) AddItem(name string, d date.Date) {
	i, found := idx.findEntry(d)
	if found {
		idx.entries[i].ItemName = name

		return
	}

	e := &DateIndexEntry{
		ItemName: name,
		Date:     d,
	}

	slices.Insert(idx.entries, i, e)
}

func (idx *DateIndex) findEntry(d date.Date) (int, bool) {
	f := func(a, b *DateIndexEntry) int {
		if a.Date.Before(b.Date) {
			return -1
		}

		if a.Date.After(b.Date) {
			return 1
		}

		return 0
	}

	dummyTgt.Date = d

	return slices.BinarySearchFunc(idx.entries, dummyTgt, f)
}

func (idx *DateIndex) FindItem(d date.Date) (string, bool) {
	i, found := idx.findEntry(d)
	if !found {
		return "", false
	}

	return idx.entries[i].ItemName, true
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

	idx.entries = entries

	if len(entries) > 0 {
		idx.first = entries[0].Date
		idx.last = sliceutils.Last(entries).Date
	}
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
