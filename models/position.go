package models

type PositionType int

const (
	PositionLong  PositionType = 1
	PositionShort PositionType = -1
)

type Position struct {
	Type           PositionType
	Entries, Exits []float64
}

func NewPosition(typ PositionType) *Position {
	return &Position{
		Type:    typ,
		Entries: []float64{},
		Exits:   []float64{},
	}
}

func (p *Position) AddEntry(currentPrice float64) {
	p.Entries = append(p.Entries, currentPrice)
}

func (p *Position) AddExit(currentPrice float64) {
	p.Exits = append(p.Exits, currentPrice)
}

func (p *Position) AnyOpen() bool {
	return len(p.Entries) > len(p.Exits)
}

func (p *Position) Liquidate(currentPrice float64) {
	for i := len(p.Exits); i < len(p.Entries); i++ {
		p.Exits = append(p.Exits, currentPrice)
	}
}

func (p *Position) OpenProfitLoss(currentPrice float64) float64 {
	var pl float64

	for i := len(p.Exits); i < len(p.Entries); i++ {
		pl += p.profitLoss(p.Entries[i], currentPrice)
	}

	return pl
}

func (p *Position) ClosedProfitLoss() float64 {
	var pl float64

	for i, exit := range p.Exits {
		pl += p.profitLoss(p.Entries[i], exit)
	}

	return pl
}

func (p *Position) profitLoss(entry, exit float64) float64 {
	if p.Type == PositionLong {
		return exit - entry
	} else if p.Type == PositionShort {
		return entry - exit
	}

	return 0.0
}
