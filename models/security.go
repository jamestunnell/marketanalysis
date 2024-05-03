package models

import (
	"fmt"
	"time"
)

const SecurityKeyName = "symbol"
const SecurityName = "security"
const SecurityNamePlural = "securities"
const SecuritySchemaStr = `{
	"$id": "https://github.com/jamestunnell/marketanalysis/security.json",
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "Security",
	"description": "Info regarding a financial security",
	"type": "object",
	"required": ["symbol", "timeZone", "open", "close"],
	"properties": {
		"symbol": {
			"type": "string",
			"minLength": 1
		},
		"timeZone": {
			"type": "string",
			"minLength": 1
		},
		"open": {"$ref": "#/$defs/timeOfDay"},
		"close": {"$ref": "#/$defs/timeOfDay"}
	},
	"$defs": {
		"timeOfDay": {
			"title": "Time-of-day",
			"description": "Hours and minutes in 24-hour format (HH:MM)",
			"type": "string",
			"pattern": "^[0-9]{2}:[0-9]{2}$"
		}
	}
}`

type Security struct {
	Symbol   string     `json:"symbol" bson:"_id"`
	TimeZone string     `json:"timeZone"`
	Open     *TimeOfDay `json:"open"`
	Close    *TimeOfDay `json:"close"`
}

func (s *Security) GetKey() string {
	return s.Symbol
}

func (s *Security) Validate() []error {
	errs := []error{}

	if s.Open.MinuteOfDay() >= s.Close.MinuteOfDay() {
		err := fmt.Errorf("open '%s' is not before close '%s'", s.Open, s.Close)

		errs = append(errs, err)
	}

	if _, err := time.LoadLocation(s.TimeZone); err != nil {
		err := fmt.Errorf("time zone '%s' is invalid: %w", s.TimeZone, err)

		errs = append(errs, err)
	}

	return errs
}
