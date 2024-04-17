package models

type BarsProvider interface {
	AnySetsLeft() bool
	Advance()
	CurrentSet() (Bars, error)
}
