package graph

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/blocks/registry"
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
	Name            string              `json:"name"`
	Type            string              `json:"type"`
	ParamVals       map[string]any      `json:"paramVals,omitempty"`
	RecordedOutputs []string            `json:"recordedOutputs,omitempty"`
	InputSources    map[string]*Address `json:"inputSources,omitempty"`
}

func (cfg Configuration) GetKey() string {
	return cfg.ID
}

func (cfg Configuration) ClearAllRecorded() {
	for _, bc := range cfg.Blocks {
		bc.RecordedOutputs = []string{}
	}
}

func (cfg Configuration) SetRecording(addr *Address) error {
	bc, found := cfg.FindBlockConfig(addr.A)
	if !found {
		return fmt.Errorf("block %s not found", addr.A)
	}

	newBlk, found := registry.Get(bc.Type)
	if !found {
		return fmt.Errorf("unknown block type %s", bc.Type)
	}

	outs := newBlk().GetOutputs()

	if _, found = outs[addr.B]; !found {
		return fmt.Errorf("block %s does not have output %s", addr.A, addr.B)
	}

	bc.RecordedOutputs = append(bc.RecordedOutputs, addr.B)

	return nil
}

func (cfg Configuration) FindBlockConfig(name string) (*BlockConfig, bool) {
	for _, bc := range cfg.Blocks {
		if bc.Name == name {
			return bc, true
		}
	}

	return nil, false
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

func (bc *BlockConfig) Validate(
	blk blocks.Block,
	findSource func(*Address) (blocks.Output, bool),
) []error {
	errs := []error{}
	ins := blk.GetInputs()

	// validate connections
	for inputName, sourceAddr := range bc.InputSources {
		in, found := ins[inputName]
		if !found {
			errs = append(errs, fmt.Errorf("block %s: input %s not found", bc.Name, inputName))

			continue
		}

		if inputName == sourceAddr.A {
			errs = append(errs, fmt.Errorf("block %s: input %s source is the same block", bc.Name, inputName))

			continue
		}

		src, found := findSource(sourceAddr)
		if !found {
			errs = append(errs, fmt.Errorf("block %s: input %s source %s not found", bc.Name, inputName, sourceAddr))
		}

		if err := src.Connect(in); err != nil {
			errs = append(errs, fmt.Errorf("block %s: cannot connect input %s to source %s: %w", bc.Name, inputName, sourceAddr, err))
		}
	}

	outs := blk.GetOutputs()

	// validate recording outputs
	for _, recOut := range bc.RecordedOutputs {
		if out, found := outs[recOut]; !found {
			errs = append(errs, fmt.Errorf("block %s: cannot reocrd output '%s' (not found)", bc.Name, recOut))
		} else if out.GetType() != "float64" {
			errs = append(errs, fmt.Errorf("block %s: cannot record output '%s' (not a float64 type)", bc.Name, recOut))
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
