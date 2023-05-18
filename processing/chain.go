package processing

import (
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type Chain struct {
	source                    Source
	procs                     []Processor
	outputProcs, outputSource float64
}

func NewChain(source Source, procs ...Processor) *Chain {
	return &Chain{
		source:       source,
		procs:        procs,
		outputProcs:  0.0,
		outputSource: 0.0,
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
	c.outputProcs = 0.0
	c.outputSource = 0.0

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

	c.outputSource = 0.0
	c.outputProcs = 0.0

	return nil
}

func (c *Chain) SourceWarm() bool {
	return c.source.Warm()
}

func (c *Chain) ProcsWarm() bool {
	return sliceutils.Last(c.procs).Warm()
}

func (c *Chain) SourceOutput() float64 {
	return c.source.Output()
}

func (c *Chain) ProcsOutput() float64 {
	return sliceutils.Last(c.procs).Output()
}

func (c *Chain) Update(bar *models.Bar) {
	c.source.Update(bar)
	if !c.source.Warm() {
		return
	}

	c.outputSource = c.source.Output()

	input := c.source.Output()

	for _, proc := range c.procs {
		proc.Update(input)

		if !proc.Warm() {
			break
		}

		input = proc.Output()
	}

}
