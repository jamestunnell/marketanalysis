package blocks

import (
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
)

type Graph struct {
	Blocks      models.Blocks
	Connections models.Connections
}

type GraphJSON struct {
	Blocks      map[string]json.RawMessage `json:"blocks"`
	Connections models.Connections         `json:"connections"`
}

func NewGraph() *Graph {
	return &Graph{
		Blocks:      models.Blocks{},
		Connections: models.Connections{},
	}
}

func (g *Graph) GetBlocks() models.Blocks {
	return g.Blocks
}

func (g *Graph) GetConnections() models.Connections {
	return g.Connections
}

func (g *Graph) MarshalJSON() ([]byte, error) {
	j := &GraphJSON{
		Blocks:      map[string]json.RawMessage{},
		Connections: g.Connections,
	}

	for name, b := range g.Blocks {
		d, err := MarshalBlockJSON(b)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to marshal block %s: %w", name, err)
		}

		j.Blocks[name] = d
	}

	return json.Marshal(j)
}

func (g *Graph) UnmarshalJSON(d []byte) error {
	var j GraphJSON

	if err := json.Unmarshal(d, &j); err != nil {
		return fmt.Errorf("failed to unmarshal as graph: %w", err)
	}

	blocks := models.Blocks{}
	for name, blockData := range j.Blocks {
		b, err := UnmarshalBlockJSON(blockData)
		if err != nil {
			return fmt.Errorf("failed to unmarshal block %s: %w", name, err)
		}

		blocks[name] = b
	}

	g.Blocks = blocks
	g.Connections = j.Connections

	return nil
}
