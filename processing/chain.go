package processing

import (
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type Chain struct {
	source       Source
	procs        []Processor
	warm         bool
	warmupPeriod int
	output       float64
}

func NewChain(source Source, procs ...Processor) *Chain {
	return &Chain{
		source:       source,
		procs:        procs,
		warm:         false,
		warmupPeriod: 0,
		output:       0.0,
	}
}

type ChainJSON struct {
	Source     json.RawMessage   `json:"source"`
	Processors []json.RawMessage `json:"processors"`
}

func (c *Chain) MarshalJSON() ([]byte, error) {
	srcJSON, err := MarshalSource(c.source)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal source JSON: %w", err)
	}

	procJSONs := make([]json.RawMessage, len(c.procs))

	for i, proc := range c.procs {
		d, err := MarshalProcessor(proc)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to marshal processor JSON: %w", err)
		}

		procJSONs[i] = d
	}

	chainJSON := &ChainJSON{
		Source:     srcJSON,
		Processors: procJSONs,
	}

	d, err := json.Marshal(chainJSON)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal chain JSON: %w", err)
	}

	return d, nil
}

func (c *Chain) UnmarshalJSON(d []byte) error {
	var chainJSON ChainJSON

	err := json.Unmarshal(d, &chainJSON)
	if err != nil {
		return fmt.Errorf("failed to unmarshal chain JSON: %w", err)
	}

	source, err := UnmarshalSource(chainJSON.Source)
	if err != nil {
		return fmt.Errorf("failed to unmarshal source JSON: %w", err)
	}

	procs := make([]Processor, len(chainJSON.Processors))
	for i, procJSON := range chainJSON.Processors {
		proc, err := UnmarshalProcessor(procJSON)
		if err != nil {
			return fmt.Errorf("failed to unmarshal processor JSON: %w", err)
		}

		procs[i] = proc
	}

	c.source = source
	c.procs = procs

	return nil
}

func (c *Chain) Initialize() error {
	if err := c.source.Initialize(); err != nil {
		return fmt.Errorf("failed to init source: %w", err)
	}

	wuPeriod := c.source.WarmupPeriod()

	for _, proc := range c.procs {
		if err := proc.Initialize(); err != nil {
			return fmt.Errorf("failed to init processor: %w", err)
		}

		wuPeriod += proc.WarmupPeriod()
	}

	c.warmupPeriod = wuPeriod
	c.warm = false
	c.output = 0.0

	return nil
}

func (c *Chain) Warm() bool {
	return c.warm
}

func (c *Chain) WarmupPeriod() int {
	return c.warmupPeriod
}

func (c *Chain) WarmUp(bars models.Bars) error {
	if len(bars) < c.warmupPeriod {
		return commonerrs.NewErrMinCount("warmup bars", len(bars), c.warmupPeriod)
	}

	err := c.source.WarmUp(bars[:c.source.WarmupPeriod()])
	if err != nil {
		return fmt.Errorf("failed to warm up source: %w", err)
	}

	if len(c.procs) == 0 {
		c.warm = true
		c.output = c.source.Output()

		return nil
	}

	wuStart := c.source.WarmupPeriod()

	for i, p := range c.procs {
		wuCount := p.WarmupPeriod()
		wuVals := make([]float64, wuCount)
		prevProcs := c.procs[:i]

		for j := 0; j < wuCount; j++ {
			b := bars[wuStart+j]

			wuVals[i] = updateChain(b, c.source, prevProcs)
		}

		err := p.WarmUp(wuVals)
		if err != nil {
			return fmt.Errorf("failed to warm up processor: %w", err)
		}
	}

	c.warm = true
	c.output = sliceutils.Last(c.procs).Output()

	return nil
}

func (c *Chain) Update(bar *models.Bar) {
	c.output = updateChain(bar, c.source, c.procs)
}

func (c *Chain) Output() float64 {
	return c.output
}

func updateChain(b *models.Bar, src Source, procs []Processor) float64 {
	src.Update(b)

	val := src.Output()

	for _, p := range procs {
		p.Update(val)

		val = p.Output()
	}

	return val
}
