package graph

import (
	"fmt"
	"strings"
)

type Address struct {
	A, B string
}

func NewAddress(a, b string) *Address {
	return &Address{
		A: a,
		B: b,
	}
}

func ParseAddress(s string) (*Address, error) {
	a := &Address{}
	if err := a.Parse(s); err != nil {
		return nil, err
	}

	return a, nil
}

func (addr *Address) String() string {
	return addr.A + "." + addr.B
}

func (addr *Address) Parse(s string) error {
	substrings := strings.Split(s, ".")
	n := len(substrings)

	if n < 2 {
		return fmt.Errorf("'%s' not formated as <A>.<B>", s)
	}

	addr.A = substrings[0]

	if n == 2 {
		addr.B = substrings[1]
	} else {
		addr.B = strings.Join(substrings[1:], ".")
	}

	return nil
}
