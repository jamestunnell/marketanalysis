package graph

import (
	"fmt"
	"slices"

	graphlib "github.com/dominikbraun/graph"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/blocks/registry"
	"github.com/jamestunnell/marketanalysis/models"
)

type Graph struct {
	*Config

	blocks       Blocks
	warmupPeriod int
	order        []string
	recordOuts   []*recordOut
	timeSeries   *models.TimeSeries
}

type recordOut struct {
	Output       blocks.Output
	Quantity     *models.Quantity
	Measurements []string
}

func New(cfg *Config) *Graph {
	log.Trace().Interface("configuration", cfg).Msg("making graph")

	return &Graph{
		Config:       cfg,
		blocks:       Blocks{},
		warmupPeriod: 0,
		order:        []string{},
		recordOuts:   []*recordOut{},
		timeSeries:   models.NewTimeSeries(),
	}
}

func (g *Graph) GetWarmupPeriod() int {
	return g.warmupPeriod
}

func (g *Graph) GetTimeSeries() *models.TimeSeries {
	return g.timeSeries
}

func (g *Graph) Init() error {
	blks, conns, err := g.makeBlocksAndConns()
	if err != nil {
		return err
	}

	if err := blks.Init(); err != nil {
		return fmt.Errorf("failed to init blocks: %w", err)
	}

	gr, err := blks.Connect(conns)
	if err != nil {
		return fmt.Errorf("failed to connect blocks: %w", err)
	}

	order, err := graphlib.TopologicalSort(gr)
	if err != nil {
		return fmt.Errorf("topological sort failed: %w", err)
	}

	log.Trace().Strs("order", order).Msg("connected graph blocks")

	wuPeriod, err := MaxTotalWarmupPeriod(blks, gr, order)
	if err != nil {
		return err
	}

	timeSeries := models.NewTimeSeries()
	recordOuts := []*recordOut{}

	// record all float64 outputs
	for blkName, blk := range blks {
		for outName, out := range blk.GetOutputs() {
			addr := &Address{A: blkName, B: outName}
			q := models.NewQuantity(addr.String())
			r := &recordOut{
				Quantity:     q,
				Output:       out,
				Measurements: g.FindMeasurements(addr),
			}

			timeSeries.AddQuantity(q)

			recordOuts = append(recordOuts, r)
		}
	}

	log.Trace().
		Int("warmupPeriod", wuPeriod).
		Msg("initialized graph model")

	g.blocks = blks
	g.warmupPeriod = wuPeriod
	g.order = order
	g.recordOuts = recordOuts
	g.timeSeries = timeSeries

	return nil
}

func (g *Graph) makeBlocksAndConns() (Blocks, []*Connection, error) {
	blks := Blocks{}
	conns := []*Connection{}

	for _, blockCfg := range g.Blocks {
		new, found := registry.Get(blockCfg.Type)
		if !found {
			err := fmt.Errorf(
				"block %s: type '%s' not found in registry",
				blockCfg.Name,
				blockCfg.Type)

			return Blocks{}, []*Connection{}, err
		}

		blk := new()
		paramVals := blockCfg.ParamVals()

		if err := blk.GetParams().SetValuesOrDefault(paramVals); err != nil {
			err = fmt.Errorf(
				"block %s: failed to set param vals %#v: %w",
				blockCfg.Name,
				paramVals,
				err)

			return Blocks{}, []*Connection{}, err
		}

		blks[blockCfg.Name] = blk

		for _, inputCfg := range blockCfg.Inputs {
			conn := &Connection{
				Source: inputCfg.Source,
				Target: &Address{A: blockCfg.Name, B: inputCfg.Name},
			}

			conns = append(conns, conn)
		}
	}

	return blks, conns, nil
}

func MaxTotalWarmupPeriod(blks Blocks, g graphlib.Graph[string, string], order []string) (int, error) {
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

		log.Trace().Str("block", name).Int("count", totalWU).Msg("total warmup")

		totalWUs[name] = totalWU
	}

	if len(totalWUs) == 0 {
		return 0, nil
	}

	return slices.Max(maps.Values(totalWUs)), nil
}

func (g *Graph) Update(bar *models.Bar, isLast bool) {
	log.Trace().Msg("updating graph")

	for _, blk := range g.blocks {
		blocks.ClearOutputs(blk)
	}

	for _, name := range g.order {
		blk := g.blocks[name]

		log.Trace().Str("name", name).Msg("running block")

		blk.Update(bar, isLast)
	}

	for _, r := range g.recordOuts {
		if r.Output.IsValueSet() {
			switch out := r.Output.(type) {
			case *blocks.TypedOutput[float64]:
				r.Quantity.AddRecord(models.NewTimeValue(bar.Timestamp, out.GetValue()))
			case *blocks.TypedOutputAsync[float64]:
				r.Quantity.AddRecord(models.NewTimeValue(out.GetTime(), out.GetValue()))
			case *blocks.TypedOutput[int]:
				r.Quantity.AddRecord(models.NewTimeValue(bar.Timestamp, float64(out.GetValue())))
			case *blocks.TypedOutputAsync[int]:
				r.Quantity.AddRecord(models.NewTimeValue(out.GetTime(), float64(out.GetValue())))
			case *blocks.TypedOutput[bool]:
				var val float64

				if out.GetValue() {
					val = 1.0
				}

				r.Quantity.AddRecord(models.NewTimeValue(bar.Timestamp, val))
			case *blocks.TypedOutputAsync[bool]:
				var val float64

				if out.GetValue() {
					val = 1.0
				}

				r.Quantity.AddRecord(models.NewTimeValue(out.GetTime(), val))
			default:
			}
		}
	}

	if !isLast {
		return
	}

	log.Trace().Msg("running last update measurements")

	// do all measurements after the last bar
	for _, r := range g.recordOuts {
		if r.Quantity.IsEmpty() {
			continue
		}

		r.Quantity.MeasureAll(r.Measurements)
	}
}
