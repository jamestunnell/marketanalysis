package models

import (
	"fmt"
	"strings"
)

type Connections map[string][]string

type Address struct {
	Block string
	Port  string
}

func ParseAddress(s string) (*Address, error) {
	substrings := strings.Split(s, ".")
	if len(substrings) != 2 {
		return nil, fmt.Errorf("'%s' not formated as <block>.<port>", s)
	}

	addr := &Address{
		Block: substrings[0],
		Port:  substrings[1],
	}

	return addr, nil
}

func (conns Connections) EachPair(each func(src, tgt *Address) error) error {
	for src, tgts := range conns {
		a, err := ParseAddress(src)
		if err != nil {
			return fmt.Errorf("invalid source address: %w", err)
		}

		for _, tgt := range tgts {
			b, err := ParseAddress(tgt)
			if err != nil {
				return fmt.Errorf("invalid target address: %w", err)
			}

			if err := each(a, b); err != nil {
				return err
			}
		}
	}

	return nil
}
