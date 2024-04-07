package models

import (
	"encoding/json"
	"fmt"
)

type Block interface {
	GetType() string
	GetDescription() string
	GetParams() Params
	GetInputs() Inputs
	GetOutputs() Outputs

	IsWarm() bool

	Init() error
	Update()
}

type Blocks map[string]Block

type NewBlockFunc func() Block

type BlockRegistry interface {
	Types() []string
	Add(typ string, newBlock NewBlockFunc)
	Get(typ string) (NewBlockFunc, bool)
}

type BlockJSON struct {
	Type   string                     `json:"type"`
	Params map[string]json.RawMessage `json:"params"`
}

func MarshalBlockJSON(blk Block) ([]byte, error) {
	ps := map[string]json.RawMessage{}
	for name, param := range blk.GetParams() {
		d, err := param.StoreVal()
		if err != nil {
			return []byte{}, fmt.Errorf("failed to store param '%s': %w", name, err)
		}

		ps[name] = d
	}

	blkJSON := &BlockJSON{
		Type:   blk.GetType(),
		Params: ps,
	}

	return json.Marshal(blkJSON)
}

func UnmarshalBlockJSON(d []byte, reg BlockRegistry) (Block, error) {
	var blk Block
	var blkJSON BlockJSON

	if err := json.Unmarshal(d, &blkJSON); err != nil {
		err = fmt.Errorf("failed to unmarshal block JSON: %w", err)

		return blk, err
	}

	newElem, found := reg.Get(blkJSON.Type)
	if !found {
		return blk, fmt.Errorf("unknown block type '%s'", blkJSON.Type)
	}

	blk = newElem()

	for name, p := range blk.GetParams() {
		rawMsg, found := blkJSON.Params[name]
		if !found {
			return blk, &ErrParamNotFound{Name: name}
		}

		err := p.LoadVal(rawMsg)
		if err != nil {
			return blk, fmt.Errorf("failed to load value for param '%s': %w", name, err)
		}
	}

	if err := blk.Init(); err != nil {
		err = fmt.Errorf("failed to init block: %w", err)

		return blk, err
	}

	return blk, nil
}
