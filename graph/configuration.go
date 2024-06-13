package graph

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/blocks/registry"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

const ConfigKeyName = "id"

type Configuration struct {
	ID     string         `json:"id" bson:"_id"`
	Name   string         `json:"name"`
	Blocks []*BlockConfig `json:"blocks"`
}

type BlockConfig struct {
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	ParamVals map[string]any  `json:"paramVals,omitempty"`
	Outputs   []*OutputConfig `json:"outputs,omitempty"`
	Inputs    []*InputConfig  `json:"inputs,omitempty"`
}

type InputConfig struct {
	Name   string   `json:"name"`
	Source *Address `json:"source"`
}

type OutputConfig struct {
	Name         string   `json:"name"`
	Measurements []string `json:"measurements"`
}

func (cfg Configuration) GetKey() string {
	return cfg.ID
}

func (cfg Configuration) FindBlock(name string) (*BlockConfig, bool) {
	for _, blk := range cfg.Blocks {
		if blk.Name == name {
			return blk, true
		}
	}

	return nil, false
}

func (cfg Configuration) FindMeasurements(outputAddr *Address) []string {
	blk, found := cfg.FindBlock(outputAddr.A)
	if !found {
		return []string{}
	}

	out, found := blk.FindOutput(outputAddr.B)
	if !found {
		return []string{}
	}

	return out.Measurements
}

func (cfg Configuration) MakeBlocks() (Blocks, []error) {
	errs := []error{}
	blks := Blocks{}

	for _, b := range cfg.Blocks {
		if _, found := blks[b.Name]; found {
			err := fmt.Errorf("duplicate block name '%s'", b.Name)

			errs = append(errs, err)

			continue
		}

		newBlk, blkFound := registry.Get(b.Type)
		if !blkFound {
			err := fmt.Errorf("block %s: has unknown type '%s'", b.Name, b.Type)

			errs = append(errs, err)

			continue
		}

		blk := newBlk()

		blks[b.Name] = blk
	}

	return blks, errs
}

func (cfg *Configuration) Validate() []error {
	blks, errs := cfg.MakeBlocks()
	if len(errs) > 0 {
		return errs
	}

	findOutput := func(addr *Address) (blocks.Output, bool) {
		blk, found := blks[addr.A]
		if !found {
			return nil, false
		}

		out, found := blk.GetOutputs()[addr.B]

		return out, found
	}

	for _, bc := range cfg.Blocks {
		if blkErrs := bc.Validate(blks[bc.Name], findOutput); len(blkErrs) > 0 {
			errs = append(errs, blkErrs...)
		}
	}

	if len(errs) == 0 {
		log.Debug().Strs("blocks", maps.Keys(blks)).Msg("blocks are all valid")
	}

	return errs
}

func (bc *BlockConfig) FindOutput(name string) (*OutputConfig, bool) {
	for _, out := range bc.Outputs {
		if out.Name == name {
			return out, true
		}
	}

	return nil, false
}

func (bc *BlockConfig) Validate(
	blk blocks.Block,
	findSource func(*Address) (blocks.Output, bool),
) []error {
	errs := []error{}
	ins := blk.GetInputs()

	// validate input config
	for _, input := range bc.Inputs {
		in, found := ins[input.Name]
		if !found {
			errs = append(errs, fmt.Errorf("block %s: input %s not found", bc.Name, input.Name))

			continue
		}

		if input.Name == input.Source.A {
			errs = append(errs, fmt.Errorf("block %s: input %s source is the same block", bc.Name, input.Name))

			continue
		}

		src, found := findSource(input.Source)
		if !found {
			errs = append(errs, fmt.Errorf("block %s: input %s source %s not found", bc.Name, input.Name, input.Source))
		}

		if err := src.Connect(in); err != nil {
			errs = append(errs, fmt.Errorf("block %s: cannot connect input %s to source %s: %w", bc.Name, input.Name, input.Source, err))
		}
	}

	outs := blk.GetOutputs()

	// validate output config
	for _, output := range bc.Outputs {
		if _, found := outs[output.Name]; !found {
			errs = append(errs, fmt.Errorf("block %s: unknown output '%s'", bc.Name, output.Name))

			continue
		}

		for _, m := range output.Measurements {
			if _, found := models.GetMeasureFunc(m); !found {
				errs = append(errs, fmt.Errorf("block %s: output %s: unknown measurement '%s'", bc.Name, output.Name, m))
			}
		}
	}

	// validate params
	params := blk.GetParams()

	for pName, val := range bc.ParamVals {
		param, found := params[pName]
		if !found {
			err := fmt.Errorf("block %s: unknown param name '%s'", bc.Name, pName)

			errs = append(errs, err)

			continue
		}

		if err := param.SetCurrentVal(val); err != nil {
			err = fmt.Errorf("block %s: param %s: value %v is invalid: %w", bc.Name, pName, val, err)

			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	if err := blk.GetParams().SetValuesOrDefault(bc.ParamVals); err != nil {
		err = fmt.Errorf("block %s: failed to set param vals %#v: %w", bc.Name, bc.ParamVals, err)

		errs = append(errs, err)
	}

	if err := blk.Init(); err != nil {
		errs = append(errs, fmt.Errorf("block %s: failed to init: %w", bc.Name, err))
	}

	return errs
}
