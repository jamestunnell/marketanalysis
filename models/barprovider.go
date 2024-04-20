package models

type BarProvider interface {
	
	AnySetsLeft() bool
	Advance()
	CurrentSet() (Bars, error)
}
