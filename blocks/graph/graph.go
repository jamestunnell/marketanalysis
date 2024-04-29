package graph

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/dominikbraun/graph"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/blocks/record"
	"github.com/jamestunnell/marketanalysis/models"
)

type Graph struct {
	Name        string
	Blocks      blocks.Blocks
	Connections blocks.Connections
	Outputs     []string

	blockOrder   []string
	recordName   string
	warmupPeriod int
}

type GraphJSON struct {
	Name        string                     `json:"name"`
	Blocks      map[string]json.RawMessage `json:"blocks"`
	Connections blocks.Connections         `json:"connections"`
	Outputs     []string                   `json:"outputs"`
}

// func LoadGraph(fpath string) (*Graph, error) {
// 	var g Graph

// 	f, err := os.Open(fpath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open model file '%s': %w", fpath, err)
// 	}

// 	decoder := json.NewDecoder(f)

// 	if err = decoder.Decode(&g); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal model JSON: %w", err)
// 	}

// 	return &g, nil
// }

func New() *Graph {
	return &Graph{
		Blocks:       blocks.Blocks{},
		Connections:  blocks.Connections{},
		blockOrder:   []string{},
		warmupPeriod: 0,
	}
}

func (m *Graph) GetName() string {
	return m.Name
}

func (m *Graph) GetBlocks() blocks.Blocks {
	return m.Blocks
}

func (m *Graph) GetConnections() blocks.Connections {
	return m.Connections
}

func (m *Graph) GetBlockOrder() []string {
	return m.blockOrder
}

func (m *Graph) GetWarmupPeriod() int {
	return m.warmupPeriod
}

// addRecorder adds a record block that uses the given recorder
// and connects outputs.
func (m *Graph) addRecorder(r blocks.Recorder) string {
	recordIns := map[string]*blocks.TypedInput[float64]{}

	for _, name := range m.Outputs {
		recordIns[name] = blocks.NewTypedInput[float64]()
	}

	record := &record.Record{
		Inputs:   recordIns,
		Recorder: r,
	}
	recordName := "record-" + nanoid.Must()

	m.Blocks[recordName] = record

	// Add connections to the record block

	for _, name := range m.Outputs {
		tgtInput := recordName + "." + name

		m.Connections[name] = append(m.Connections[name], tgtInput)
	}

	return recordName
}

func (m *Graph) Init(rec blocks.Recorder) error {
	if m.recordName != "" {
		delete(m.Blocks, m.recordName)
	}

	m.recordName = m.addRecorder(rec)

	if err := m.Blocks.Init(); err != nil {
		return fmt.Errorf("failed to init blocks: %w", err)
	}

	g, err := m.Blocks.Connect(m.Connections)
	if err != nil {
		return fmt.Errorf("failed to connect blocks: %w", err)
	}

	order, err := graph.TopologicalSort(g)
	if err != nil {
		return fmt.Errorf("topological sort failed: %w", err)
	}

	log.Debug().Strs("order", order).Msg("connected graph blocks")

	totalWUs := map[string]int{}

	predMap, err := g.PredecessorMap()
	if err != nil {
		return fmt.Errorf("failed to make predecessor map: %w", err)
	}

	for _, name := range order {
		predTotalWUs := []int{}
		for predName := range predMap[name] {
			predTotalWUs = append(predTotalWUs, totalWUs[predName])
		}

		totalWU := m.Blocks[name].GetWarmupPeriod()
		if len(predTotalWUs) > 0 {
			totalWU += slices.Max(predTotalWUs)
		}

		log.Debug().Str("block", name).Int("count", totalWU).Msg("total warmup")

		totalWUs[name] = totalWU
	}

	m.blockOrder = order
	m.warmupPeriod = slices.Max(maps.Values(totalWUs))

	log.Debug().
		Int("warmupPeriod", m.warmupPeriod).
		Msg("initialized graph model")

	return nil
}

func (m *Graph) Update(bar *models.Bar) {
	for _, blockName := range m.blockOrder {
		blk, found := m.Blocks[blockName]
		if !found {
			log.Fatal().Str("name", blockName).Msg("block not found")
		}

		blk.Update(bar)
	}
}

// func (m *Graph) MarshalJSON() ([]byte, error) {
// 	j := &GraphJSON{
// 		Blocks:      map[string]json.RawMessage{},
// 		Connections: m.Connections,
// 		Outputs:     m.Outputs,
// 	}

// 	for name, b := range m.Blocks {
// 		d, err := MarshalBlockJSON(b)
// 		if err != nil {
// 			return []byte{}, fmt.Errorf("failed to marshal block %s: %w", name, err)
// 		}

// 		j.Blocks[name] = d
// 	}

// 	return json.Marshal(j)
// }

// func (m *Graph) UnmarshalJSON(d []byte) error {
// 	var j GraphJSON

// 	if err := json.Unmarshal(d, &j); err != nil {
// 		return fmt.Errorf("failed to unmarshal as graph: %w", err)
// 	}

// 	blocks := blocks.Blocks{}
// 	for name, blockData := range j.Blocks {
// 		b, err := UnmarshalBlockJSON(blockData)
// 		if err != nil {
// 			return fmt.Errorf("failed to unmarshal block %s: %w", name, err)
// 		}

// 		blocks[name] = b
// 	}

// 	m.Blocks = blocks
// 	m.Connections = j.Connections
// 	m.Outputs = j.Outputs

// 	return nil
// }
