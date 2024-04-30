package graph

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/jamestunnell/marketanalysis/blocks"
// 	"github.com/jamestunnell/marketanalysis/blocks.Registry"
// 	"github.com/jamestunnell/marketanalysis/commonerrs"
// )

// type BlockJSON struct {
// 	Type   string                     `json:"type"`
// 	Params map[string]json.RawMessage `json:"params"`
// }

// func MarshalBlockJSON(blk blocks.Block) ([]byte, error) {
// 	ps := map[string]json.RawMessage{}
// 	for name, param := range blk.GetParams() {
// 		d, err := param.StoreVal()
// 		if err != nil {
// 			return []byte{}, fmt.Errorf("failed to store param '%s': %w", name, err)
// 		}

// 		ps[name] = d
// 	}

// 	blkJSON := &BlockJSON{
// 		Type:   blk.GetType(),
// 		Params: ps,
// 	}

// 	return json.Marshal(blkJSON)
// }

// func UnmarshalBlockJSON(d []byte) (blocks.Block, error) {
// 	var blk blocks.Block
// 	var blkJSON BlockJSON

// 	if err := json.Unmarshal(d, &blkJSON); err != nil {
// 		err = fmt.Errorf("failed to unmarshal block JSON: %w", err)

// 		return blk, err
// 	}

// 	entry, found := registry.Get(blkJSON.Type)
// 	if !found {
// 		return blk, fmt.Errorf("unknown block type '%s'", blkJSON.Type)
// 	}

// 	blk = entry.New()

// 	for name, p := range blk.GetParams() {
// 		rawMsg, found := blkJSON.Params[name]
// 		if !found {
// 			return nil, commonerrs.NewErrNotFound("param", name)
// 		}

// 		if err := p.LoadVal(rawMsg); err != nil {
// 			return nil, fmt.Errorf("failed to load value for param '%s': %w", name, err)
// 		}
// 	}

// 	return blk, nil
// }
