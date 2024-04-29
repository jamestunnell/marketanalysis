package models

const GraphDefKeyName = "id"
const GraphDefName = "graph"
const GraphDefNamePlural = "graphs"
const GraphDefSchemaStr = `{
	"$id": "https://github.com/jamestunnell/marketanalysis/model.json",
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "Graph Definition",
	"description": "Define a graph of connected blocks",
	"type": "object",
	"required": ["id", "name", "blocks", "connections"],
	"properties": {
		"id": {
			"type": "string",
			"minLength": 1
		},
		"name": {
			"type": "string",
			"minLength": 1
		},
		"blocks": {
			"title": "Blocks",
			"description": "Blocks to include in the graph",
			"type": "array",
			"items": {"$ref": "#/$defs/blockUse"}
		},
		"connections": {
			"type": "array",
			"items": {"$ref": "#/$defs/connection"}
		}
	},
	"$defs": {
		"blockUse": {
			"type": "object",
			"required": ["name", "type"],
			"properties": {
				"name": {
					"type": "string",
					"minLength": 1
				},
				"type": {
					"type": "string",
					"minLength": 1
				},
				"recording": {
					"title": "Recording",
					"descriptions": "Outputs to record",
					"type": "array",
					"items": {
						"type": "string"
					}
				}
			}
		},
		"connection": {
			"type": "object",
			"required": ["source", "targets"],
			"properties": {
				"source": {"$ref": "#/$defs/portAddr"},
				"targets": {"$ref": "#/$defs/portAddrArray"}
			}
		},
		"portAddrArray": {
			"type": "array",
			"items": {"$ref": "#/$defs/portAddr"}
		},
		"portAddr": {
			"type": "string",
			"minLength": 3,
			"pattern": "^[_A-Za-z]+.[_A-Za-z]+$"
		}
	}
}`

type GraphDef struct {
	ID          string        `json:"id" bson:"_id"`
	Name        string        `json:"name"`
	Blocks      []*BlockUse   `json:"blocks"`
	Connections []*Connection `json:"connections"`
	Outputs     []string      `json:"outputs"`
}

type BlockUse struct {
	*NameType

	ParamVals map[string]any `json:"paramVals"`
	Recording []string       `json:"recording"`
}

type NameType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Connection struct {
	Source  string   `json:"source"`
	Targets []string `json:"targets"`
}
