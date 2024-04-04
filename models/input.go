package models

type Input interface {
	Name() string
	Type() string
}

type Inputs map[string]Input
