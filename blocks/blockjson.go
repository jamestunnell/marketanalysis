package blocks

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/jamestunnell/marketanalysis/models"
)

type BlockJSON struct {
	Type   string                     `json:"type"`
	Params map[string]json.RawMessage `json:"params"`
}

type errMissingParam struct {
	Name string
}

type errInvalidParam struct {
	Name   string
	Errors []error
}

func MarshalBlockJSON(blk models.Block) ([]byte, error) {
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

func UnmarshalBlockJSON(d []byte) (models.Block, error) {
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
			return nil, &errMissingParam{Name: name}
		}

		if err := p.LoadVal(rawMsg); err != nil {
			return nil, fmt.Errorf("failed to load value for param '%s': %w", name, err)
		}

		if errs := models.ValidateParam(p); len(errs) > 0 {
			return nil, &errInvalidParam{Name: name, Errors: errs}
		}
	}

	return blk, nil
}

func (err *errMissingParam) Error() string {
	return fmt.Sprintf("missing param %s", err.Name)
}

func (err *errInvalidParam) Error() string {
	var merr *multierror.Error

	for _, err := range err.Errors {
		merr = multierror.Append(merr, err)
	}

	return fmt.Sprintf("invalid param %s: %v", err.Name, merr)
}
