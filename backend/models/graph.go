package models

const GraphKeyName = "id"
const GraphName = "graph"
const GraphNamePlural = "graphs"
const GraphSchemaStr = `{
	"$id": "https://github.com/jamestunnell/marketanalysis/model.json",
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "Graph",
	"description": "Graph of blocks",
	"type": "object",
	"required": ["id", "name", "blocks", "connections", "outputs"],
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
			"type": "array",
			"items": {"$ref": "#/$defs/graphBlock"}
		},
		"connections": {
			"type": "array",
			"items": {"$ref": "#/$defs/graphConnection"}
		},
		"outputs": {"$ref": "#/$defs/portAddrArray"}
	},
	"$defs": {
		"graphBlock": {
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
				}
			}
		},
		"graphConnection": {
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

type Graph struct {
	ID          string             `json:"id" bson:"_id"`
	Name        string             `json:"name"`
	Blocks      []*GraphBlock      `json:"blocks"`
	Connections []*GraphConnection `json:"connections"`
	Outputs     []string           `json:"outputs"`
}

type GraphBlock struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type GraphConnection struct {
	Source  string   `json:"source"`
	Targets []string `json:"targets"`
}
