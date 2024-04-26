package models

import (
	m "github.com/jamestunnell/marketanalysis/models"
)

const SecurityKeyName = "symbol"
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
	Symbol   string       `json:"symbol" bson:"_id"`
	TimeZone string       `json:"timeZone"`
	Open     *m.TimeOfDay `json:"open"`
	Close    *m.TimeOfDay `json:"close"`
}
