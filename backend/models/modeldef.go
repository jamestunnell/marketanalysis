package models

const ModelDefKeyName = "id"
const ModelDefNamePlural = "modeldefs"
const ModelDefSchemaStr = `{
	"$id": "https://github.com/jamestunnell/marketanalysis/modeldef.json",
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "Model definition",
	"description": "Define a financial model",
	"type": "object",
	"required": ["id", "name", "type", "paramDefs"],
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
		"paramDefs": {"$ref": "#/$defs/paramDef"}
	},
	"$defs": {
		"paramDef": {
			"title": "Parameter definition",
			"type": "object",
			"required": ["name", "type"],
		}
	}
}`

type ModelDef struct {
	ID     string               `json:"id" bson:"_id"`
	Name   string               `json:"name"`
	Type   string               `json:"type"`
	Params map[string]*ParamDef `json:"paramDefs"`
}

type ParamDef struct {
}
