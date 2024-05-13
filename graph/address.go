package graph

import (
	"encoding/json"
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

func (addr *Address) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(addr.String())
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal JSON as string: %w", err)
	}

	return d, nil
}

func (addr *Address) UnmarshalJSON(d []byte) error {
	var str string

	if err := json.Unmarshal(d, &str); err != nil {
		return fmt.Errorf("failed to unmarshal JSON as string: %w", err)
	}

	if err := addr.Parse(str); err != nil {
		return fmt.Errorf("failed to parse '%s' as address: %w", str, err)
	}

	return nil
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

	for _, sub := range substrings {
		if sub == "" {
			return fmt.Errorf("address contains empty substring")
		}
	}

	addr.A = substrings[0]

	if n == 2 {
		addr.B = substrings[1]
	} else {
		addr.B = strings.Join(substrings[1:], ".")
	}

	return nil
}
