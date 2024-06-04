package models

import (
	"fmt"
	"time"
)

type Security struct {
	Symbol   string `json:"symbol" bson:"_id"`
	Days     int    `json:"days"`
	TimeZone string `json:"timeZone"`
}

const SecurityKeyName = "symbol"

func (s *Security) GetKey() string {
	return s.Symbol
}

func (s *Security) Validate() []error {
	errs := []error{}

	_, err := time.LoadLocation(s.TimeZone)
	if err != nil {
		err = fmt.Errorf("time zone '%s' is invalid: %w", s.TimeZone, err)

		errs = append(errs, err)
	}

	if s.Days < 1 {
		err = fmt.Errorf("days '%d' must be > 1", s.Days)

		errs = append(errs, err)
	}

	return errs
}
