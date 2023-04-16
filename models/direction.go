package models

type Direction int

const (
	DirNone Direction = iota
	DirDown
	DirUp
)

func (d Direction) String() string {
	var str string

	switch d {
	case DirNone:
		str = "none"
	case DirDown:
		str = "down"
	case DirUp:
		str = "up"
	}

	return str
}
