package models

import "github.com/rickb777/date"

const RunDay = "Day"

type RunDayRequest struct {
	Symbol  string    `json:"symbol"`
	Date    date.Date `json:"date"`
	Format  string    `json:"format"`
	LocalTZ string    `json:"localTZ,omitempty"`
}
