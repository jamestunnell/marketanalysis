package models

import (
	"encoding/json"
	"fmt"
)

type Graph struct {
	Blocks      Blocks
	Connections []*Connection
}

type GraphJSON struct {
	Blocks      map[string]json.RawMessage `json:"blocks"`
	Connections []*Connection              `json:"connections"`
}

func NewGraph() *Graph {
	return &Graph{
		Blocks:      Blocks{},
		Connections: []*Connection{},
	}
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

func UnmarshalGraphJSON(d []byte, reg BlockRegistry) (*Graph, error) {
	var j GraphJSON

	if err := json.Unmarshal(d, &j); err != nil {
		return nil, fmt.Errorf("failed to unmarshal as graph: %w", err)
	}

	blocks := Blocks{}
	for name, blockData := range j.Blocks {
		b, err := UnmarshalBlockJSON(blockData, reg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal block %s: %w", name, err)
		}

		blocks[name] = b
	}

	g := &Graph{
		Blocks:      blocks,
		Connections: j.Connections,
	}

	return g, nil
}
