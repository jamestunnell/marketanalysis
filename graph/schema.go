package graph

import (
	"github.com/rs/zerolog/log"

	"github.com/xeipuuv/gojsonschema"
)

var configSchema *gojsonschema.Schema

func init() {
	l := gojsonschema.NewStringLoader(ConfigSchemaStr)

	schema, err := gojsonschema.NewSchema(l)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make graph config schema")
	}

	configSchema = schema
}

const ConfigSchemaStr = `{
	"$id": "https://github.com/jamestunnell/marketanalysis/graph/config.json",
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "Graph configuration",
	"description": "A graph of connected blocks",
	"type": "object",
	"required": ["id", "blocks", "connections"],
	"properties": {
		"id": {
			"type": "string",
			"pattern": "^[A-Za-z0-9_\\-]+$"
		},
		"name": { "type": "string" },
		"blocks": {
			"title": "Blocks",
			"description": "Map block names to config",
			"type": "object",
			"patternProperties": {
				"^[A-Za-z_]+$": {"$ref": "#/$defs/blockConfig"}
			},
			"additionalProperties": false
		},
		"connections": {
			"type": "array",
			"items": {"$ref": "#/$defs/connection"}
		}
	},
	"$defs": {
		"blockConfig": {
			"title": "Block configuration",
			"description": "Config needed to make a block",
			"type": "object",
			"required": ["type"],
			"properties": {
				"type": {
					"type": "string",
					"minLength": 1
				},
				"paramVals": {
					"title": "Param values",
					"description": "Non-default values to use",
					"type": "object",
					"patternProperties": {
						"^[A-Za-z_]+$": {
							"type": ["boolean", "integer", "number", "string"]
						}
					},
					"additionalProperties": false
				},
				"recording": {
					"title": "Recording",
					"descriptions": "Outputs to record",
					"type": "array",
					"items": {
						"type": "string",
						"minLength": 1
					}
				}
			}
		},
		"connection": {
			"type": "object",
			"required": ["source", "target"],
			"properties": {
				"source": {"$ref": "#/$defs/addr"},
				"target": {"$ref": "#/$defs/addr"}
			}
		},
		"addr": {
			"type": "string",
			"minLength": 3,
			"pattern": "^[_A-Za-z]+.[_A-Za-z]+$"
		}
	}
}`

func GetConfigSchema() *gojsonschema.Schema {
	return configSchema
}
