package graph

import (
	"fmt"
	"slices"

	"github.com/dominikbraun/graph"
	gr "github.com/dominikbraun/graph"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/blocks/record"
	"github.com/jamestunnell/marketanalysis/blocks/registry"
	"github.com/jamestunnell/marketanalysis/models"
)

type Graph struct {
	*Configuration

	blocks       Blocks
	warmupPeriod int
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

func New(cfg *Configuration) *Graph {
	return &Graph{
		Configuration: cfg,
		blocks:        Blocks{},
		warmupPeriod:  0,
	}
}

func (m *Graph) GetWarmupPeriod() int {
	return m.warmupPeriod
}

func (m *Graph) Init(rec blocks.Recorder) error {
	blks, conns, err := m.makeBlocksAndConns(rec)
	if err != nil {
		return err
	}

	if err := blks.Init(); err != nil {
		return fmt.Errorf("failed to init blocks: %w", err)
	}

	g, err := blks.Connect(conns)
	if err != nil {
		return fmt.Errorf("failed to connect blocks: %w", err)
	}

	order, err := graph.TopologicalSort(g)
	if err != nil {
		return fmt.Errorf("topological sort failed: %w", err)
	}

	log.Debug().Strs("order", order).Msg("connected graph blocks")

	wuPeriod, err := MaxTotalWarmupPeriod(blks, g, order)
	if err != nil {
		return err
	}

	log.Debug().
		Int("warmupPeriod", wuPeriod).
		Msg("initialized graph model")

	m.blocks = blks
	m.warmupPeriod = wuPeriod

	return nil
}

func (m *Graph) makeBlocksAndConns(r blocks.Recorder) (Blocks, []*Connection, error) {
	blks := Blocks{}
	conns := slices.Clone(m.Connections)
	recordName := "record-" + nanoid.Must()
	recordIns := map[string]*blocks.TypedInput[float64]{}

	for _, cfg := range m.Blocks {
		new, found := registry.Get(cfg.Type)
		if !found {
			err := fmt.Errorf("block %s: type '%s' not found in registry", cfg.Name, cfg.Type)

			return Blocks{}, []*Connection{}, err
		}

		blk := new()

		if err := blk.GetParams().SetValuesOrDefault(cfg.ParamVals); err != nil {
			err = fmt.Errorf("block %s: failed to set param vals %#v: %w", cfg.Name, cfg.ParamVals, err)

			return Blocks{}, []*Connection{}, err
		}

		blks[cfg.Name] = blk

		for _, outName := range cfg.Recording {
			if _, found := blk.GetOutputs()[outName]; !found {
				err := fmt.Errorf("block %s: recording output '%s' not found", cfg.Name, outName)

				return Blocks{}, []*Connection{}, err
			}

			recTarget := fmt.Sprintf("%s.%s", cfg.Name, outName)
			recordConn := &Connection{
				Source: NewAddress(cfg.Name, outName),
				Target: NewAddress(recordName, recTarget),
			}

			recordIns[recTarget] = blocks.NewTypedInput[float64]()

			conns = append(conns, recordConn)
		}
	}

	blks[recordName] = &record.Record{
		Inputs:   recordIns,
		Recorder: r,
	}

	return blks, conns, nil
}

func MaxTotalWarmupPeriod(blks Blocks, g gr.Graph[string, string], order []string) (int, error) {
	totalWUs := map[string]int{}

	predMap, err := g.PredecessorMap()
	if err != nil {
		return 0, fmt.Errorf("failed to make predecessor map: %w", err)
	}

	for _, name := range order {
		predTotalWUs := []int{}
		for predName := range predMap[name] {
			predTotalWUs = append(predTotalWUs, totalWUs[predName])
		}

		totalWU := blks[name].GetWarmupPeriod()
		if len(predTotalWUs) > 0 {
			totalWU += slices.Max(predTotalWUs)
		}

		log.Debug().Str("block", name).Int("count", totalWU).Msg("total warmup")

		totalWUs[name] = totalWU
	}

	maxWU := slices.Max(maps.Values(totalWUs))

	return maxWU, nil
}

func (m *Graph) Update(bar *models.Bar) {
	for _, blk := range m.blocks {
		blk.Update(bar)
	}
}
