package models

type Output interface {
	Name() string
	Type() string

	Connect(Input)
}

type Outputs map[string]Output
