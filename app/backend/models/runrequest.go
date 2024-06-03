package models

import "github.com/rickb777/date"

const RunDay = "day"

type RunDayRequest struct {
	Symbol     string    `json:"symbol"`
	Date       date.Date `json:"date"`
	Format     string    `json:"format"`
	TimeZone   string    `json:"timeZone"`
	ShowWarmup bool      `json:"showWarmup,omitempty"`
}
