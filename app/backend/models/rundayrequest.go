package models

import "github.com/rickb777/date"

type RunDayRequest struct {
	Symbol  string    `json:"symbol"`
	Date    date.Date `json:"date"`
	LocalTZ string    `json:"localTZ,omitempty"`
}
