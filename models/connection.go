package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Connection struct {
	Source *Address `json:"source"`
	Dest   *Address `json:"dest"`
}

type Address struct {
	Block string
	Port  string
}

func NewAddress(block, port string) *Address {
	return &Address{
		Block: block,
		Port:  port,
	}
}

func NewConnection(src, dest *Address) *Connection {
	return &Connection{
		Source: src,
		Dest:   dest,
	}
}

func (addr *Address) String() string {
	return addr.Block + "." + addr.Port
}

func (addr *Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(addr.String())
}

func (addr *Address) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("failed to unmarshal address JSON as string: %w", err)
	}

	return addr.Parse(str)
}

func (addr *Address) Parse(s string) error {
	substrings := strings.Split(s, ".")
	if len(substrings) != 2 {
		return fmt.Errorf("string '%s' not formated as <block>.<port>", s)
	}

	addr.Block = substrings[0]
	addr.Port = substrings[1]

	return nil
}
