package graph

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"

	"github.com/jamestunnell/marketanalysis/blocks/registry"
)

const ConfigKeyName = "id"

type Configuration struct {
	ID          string         `json:"id" bson:"_id"`
	Name        string         `json:"name"`
	Blocks      []*BlockConfig `json:"blocks"`
	Connections []*Connection  `json:"connections"`
}

type BlockConfig struct {
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	ParamVals map[string]any `json:"paramVals"`
	Recording []string       `json:"recording"`
}

func (cfg Configuration) GetKey() string {
	return cfg.ID
}

func (cfg Configuration) ClearAllRecording() {
	for _, bc := range cfg.Blocks {
		bc.Recording = []string{}
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

	bc.Recording = append(bc.Recording, addr.B)

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

	for bName, blk := range blks {
		bc, _ := cfg.FindBlockConfig(bName)
		outs := blk.GetOutputs()

		// validate recording outputs
		for _, recOut := range bc.Recording {
			if out, found := outs[recOut]; !found {
				errs = append(errs, fmt.Errorf("block %s: recording output '%s' not found", bName, recOut))
			} else if out.GetType() != "float64" {
				errs = append(errs, fmt.Errorf("block %s: recording output '%s' not a float64 type", bName, recOut))
			}
		}

		// validate params
		params := blk.GetParams()

		for pName, val := range bc.ParamVals {
			param, found := params[pName]
			if !found {
				err := fmt.Errorf("block %s: unknown param name '%s'", bName, pName)

				errs = append(errs, err)

				continue
			}

			if err := param.SetVal(val); err != nil {
				err = fmt.Errorf("block %s: param %s: value %v is invalid: %w", bName, pName, val, err)

				errs = append(errs, err)
			}
		}

		if err := blk.GetParams().SetValuesOrDefault(bc.ParamVals); err != nil {
			err = fmt.Errorf("block %s: failed to set param vals %#v: %w", bName, bc.ParamVals, err)

			errs = append(errs, err)
		}

		if err := blk.Init(); err != nil {
			errs = append(errs, fmt.Errorf("block %s: failed to init: %w", bName, err))
		}
	}

	if len(errs) == 0 {
		log.Debug().Strs("blocks", maps.Keys(blks)).Msg("blocks are all valid")
	}

	// validate connections
	for i, conn := range cfg.Connections {
		if conn == nil {
			err := fmt.Errorf("connection #%d is null", i+1)

			errs = append(errs, err)

			continue
		}

		if _, found := blks.FindOutput(conn.Source); !found {
			err := fmt.Errorf("connection #%d: source output %s not found", i+1, conn.Source)

			errs = append(errs, err)
		}

		if _, found := blks.FindInput(conn.Target); !found {
			err := fmt.Errorf("connection #%d: target input %s not found", i+1, conn.Target)

			errs = append(errs, err)
		}
	}

	return errs
}
