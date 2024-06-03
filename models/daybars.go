package models

import (
	"errors"
	"slices"

	"github.com/rickb777/date"
)

type DayBars struct {
	Date date.Date `json:"date" bson:"_id"`
	Bars Bars      `json:"bars"`
}

var (
	errBarsNotSorted         = errors.New("bars are not sorted")
	errBarWithUnexpectedDate = errors.New("bar with unexpected date")
)

func (db *DayBars) GetKey() string {
	return db.Date.String()
}

func (db *DayBars) Validate() []error {
	errs := []error{}

	if !slices.IsSortedFunc(db.Bars, CompareBarsByTimestamp) {
		err := errBarsNotSorted

		errs = append(errs, err)
	}

	for _, bar := range db.Bars {
		if !date.NewAt(bar.Timestamp).Equal(db.Date) {
			errs = append(errs, errBarWithUnexpectedDate)

			break
		}
	}

	return errs
}

func CompareBarsByTimestamp(a, b *Bar) int {
	return a.Timestamp.Compare(b.Timestamp)
}
