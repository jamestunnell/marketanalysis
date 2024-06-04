package models

import (
	"errors"
	"slices"

	"github.com/rickb777/date"
)

type BarSet struct {
	Date string `json:"date" bson:"_id"`
	Bars Bars   `json:"bars"`
}

var (
	errBarsNotSorted         = errors.New("bars are not sorted")
	errBarWithUnexpectedDate = errors.New("bar with unexpected date")
)

func (db *BarSet) GetKey() string {
	return db.Date
}

func (db *BarSet) Validate() []error {
	errs := []error{}

	if !slices.IsSortedFunc(db.Bars, CompareBarsByTimestamp) {
		err := errBarsNotSorted

		errs = append(errs, err)
	}

	for _, bar := range db.Bars {
		if date.NewAt(bar.Timestamp).String() != db.Date {
			errs = append(errs, errBarWithUnexpectedDate)

			break
		}
	}

	return errs
}

func CompareBarsByTimestamp(a, b *Bar) int {
	return a.Timestamp.Compare(b.Timestamp)
}
