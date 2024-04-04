package blocks

import (
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
)

type BlockJSON struct {
	Type   string                     `json:"type"`
	Params map[string]json.RawMessage `json:"params"`
}

type errMissingParam struct {
	Name string
}

func MarshalJSON(blk models.Block) ([]byte, error) {
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

func UnmarshalJSON(d []byte) (models.Block, error) {
	var blk models.Block
	var blkJSON BlockJSON

	if err := json.Unmarshal(d, &blkJSON); err != nil {
		err = fmt.Errorf("failed to unmarshal block JSON: %w", err)

		return blk, err
	}

	newElem, found := Registry().Get(blkJSON.Type)
	if !found {
		return blk, fmt.Errorf("unknown block type '%s'", blkJSON.Type)
	}

	blk = newElem()

	for name, p := range blk.GetParams() {
		rawMsg, found := blkJSON.Params[name]
		if !found {
			return blk, &errMissingParam{Name: name}
		}

		err := p.LoadVal(rawMsg)
		if err != nil {
			return blk, fmt.Errorf("failed to load value for param '%s': %w", name, err)
		}
	}

	if err := blk.Init(); err != nil {
		err = fmt.Errorf("failed to init blkent: %w", err)

		return blk, err
	}

	return blk, nil
}

func (err *errMissingParam) Error() string {
	return fmt.Sprintf("missing param %s", err.Name)
}