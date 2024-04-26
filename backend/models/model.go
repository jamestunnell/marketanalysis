package models

const ModelKeyName = "id"
const ModelNamePlural = "models"
const ModelSchemaStr = `{
	"$id": "https://github.com/jamestunnell/marketanalysis/model.json",
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "Model definition",
	"description": "Define a financial model",
	"type": "object",
	"required": ["id", "name", "type", "params"],
	"properties": {
		"id": {
			"type": "string",
			"minLength": 1
		},
		"name": {
			"type": "string",
			"minLength": 1
		},
		"type": {
			"type": "string",
			"minLength": 1
		},
		"params": {
			"type": "array",
			"items": {"$ref": "#/$defs/param"}
		}
	},
	"$defs": {
		"param": {
			"title": "Parameter definition",
			"type": "object",
			"required": ["name", "schema"],
			"properties": {
				"name": {
					"type": "string",
					"minLength": 1
				},
				"schema": {
					"type": "object",
				},
		}
	}
}`

type Model struct {
	ID     string            `json:"id" bson:"_id"`
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Params map[string]*Param `json:"params"`
}

type Param struct {
	Name   string         `json:"name"`
	Schema map[string]any `json:"schema"`
}
