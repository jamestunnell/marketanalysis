package blocks

import (
	"fmt"
)

// Connections maps source address string to target address strings.
type Connections map[string][]string

func (conns Connections) EachPair(each func(src, tgt *Address) error) error {
	for srcString, tgtStrings := range conns {
		src := &Address{}

		if err := src.Parse(srcString); err != nil {
			return fmt.Errorf("invalid source address: %w", err)
		}

		for _, tgtStr := range tgtStrings {
			tgt := &Address{}

			err := tgt.Parse(tgtStr)
			if err != nil {
				return fmt.Errorf("invalid target address: %w", err)
			}

			if err := each(src, tgt); err != nil {
				return err
			}
		}
	}

	return nil
}
