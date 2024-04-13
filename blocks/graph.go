package blocks

import (
	"encoding/json"
	"fmt"
	"os"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/models"
)

type Graph struct {
	Name        string
	Blocks      models.Blocks
	Connections models.Connections
	Outputs     []string

	blockOrder []string
}

type GraphJSON struct {
	Name        string                     `json:"name"`
	Blocks      map[string]json.RawMessage `json:"blocks"`
	Connections models.Connections         `json:"connections"`
	Outputs     []string                   `json:"outputs"`
}

const TypeGraph = "Graph"

func LoadGraphFile(fpath string) (*Graph, error) {
	var g Graph

	f, err := os.Open(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed to open model file '%s': %w", fpath, err)
	}

	decoder := json.NewDecoder(f)

	if err = decoder.Decode(&g); err != nil {
		return nil, fmt.Errorf("failed to unmarshal model JSON: %w", err)
	}

	return &g, nil
}

func NewGraph() *Graph {
	return &Graph{
		Blocks:      models.Blocks{},
		Connections: models.Connections{},
		blockOrder:  []string{},
	}
}

func (g *Graph) GetType() string {
	return TypeGraph
}

func (g *Graph) GetName() string {
	return g.Name
}

func (g *Graph) GetBlocks() models.Blocks {
	return g.Blocks
}

func (g *Graph) GetConnections() models.Connections {
	return g.Connections
}

func (g *Graph) GetBlockOrder() []string {
	return g.blockOrder
}

func (m *Graph) GetParams() models.Params {
	params := models.Params{}

	for blkName, blk := range m.Blocks {
		for paramName, param := range blk.GetParams() {
			combinedName := models.NewAddress(blkName, paramName).String()

			params[combinedName] = param
		}
	}

	return params
}

func (m *Graph) GetOutputs() models.Outputs {
	outs := models.Outputs{}

	for blkName, blk := range m.Blocks {
		for outName, out := range blk.GetOutputs() {
			combinedName := models.NewAddress(blkName, outName).String()

			outs[combinedName] = out
		}
	}

	return outs
}

func (g *Graph) AddRecorder(r models.Recorder) {
	// Add the record block

	recordIns := map[string]*models.TypedInput[float64]{}

	for _, name := range g.Outputs {
		recordIns[name] = models.NewTypedInput[float64]()
	}

	record := &Record{
		Inputs:   recordIns,
		Recorder: r,
	}
	recordName := "record-" + nanoid.Must()

	g.Blocks[recordName] = record

	// Add connections to the record block

	for _, name := range g.Outputs {
		tgtInput := recordName + "." + name

		g.Connections[name] = append(g.Connections[name], tgtInput)
	}
}

func (g *Graph) Init() error {
	if err := g.Blocks.Init(); err != nil {
		return fmt.Errorf("failed to init blocks: %w", err)
	}

	order, err := g.Blocks.Connect(g.Connections)
	if err != nil {
		return fmt.Errorf("failed to connect blocks: %w", err)
	}

	log.Debug().Strs("order", order).Msg("connected graph blocks")

	g.blockOrder = order

	return nil
}

func (g *Graph) Update(bar *models.Bar) {
	// reset all inputs so we know what has been set this cycle
	for _, blk := range g.Blocks {
		for _, in := range blk.GetInputs() {
			in.Reset()
		}
	}

	for _, blockName := range g.blockOrder {
		blk, found := g.Blocks[blockName]
		if !found {
			log.Fatal().Str("name", blockName).Msg("block not found")
		}

		blk.Update(bar)
	}
}

func (g *Graph) MarshalJSON() ([]byte, error) {
	j := &GraphJSON{
		Blocks:      map[string]json.RawMessage{},
		Connections: g.Connections,
		Outputs:     g.Outputs,
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
	g.Outputs = j.Outputs

	return nil
}
