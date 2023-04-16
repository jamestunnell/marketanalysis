package models

import "time"

type PositionType int

const (
	PositionTypeLong  = "long"
	PositionTypeShort = "short"
)

type Position struct {
	Type       string     `json:"type"`
	Entry      *TimePrice `json:"entry"`
	Exit       *TimePrice `json:"exit"`
	ExitReason string     `json:"exitReason"`
	ClosedPL   float64    `json:"closedPL"`
}

func NewLongPosition(t time.Time, price float64) *Position {
	entry := &TimePrice{
		Time:  t,
		Price: price,
	}

	return &Position{
		Type:       PositionTypeLong,
		Entry:      entry,
		Exit:       nil,
		ExitReason: "",
		ClosedPL:   0.0,
	}
}

func NewShortPosition(t time.Time, price float64) *Position {
	entry := &TimePrice{
		Time:  t,
		Price: price,
	}

	return &Position{
		Type:       PositionTypeShort,
		Entry:      entry,
		Exit:       nil,
		ExitReason: "",
		ClosedPL:   0.0,
	}
}

func (p *Position) Localize() {
	if p.Entry != nil {
		p.Entry.Time = p.Entry.Time.Local()
	}

	if p.Exit != nil {
		p.Exit.Time = p.Exit.Time.Local()
	}
}

func (p *Position) IsOpen() bool {
	return p.Exit == nil
}

func (p *Position) Close(t time.Time, price float64, reason string) {
	p.Exit = &TimePrice{Time: t, Price: price}
	p.ExitReason = reason

	if p.Type == PositionTypeLong {
		p.ClosedPL = p.Exit.Price - p.Entry.Price
	} else if p.Type == PositionTypeShort {
		p.ClosedPL = p.Entry.Price - p.Exit.Price
	}
}

func (p *Position) OpenProfitLoss(currentPrice float64) (float64, bool) {
	if !p.IsOpen() {
		return 0.0, false
	}

	if p.Type == PositionTypeLong {
		return currentPrice - p.Entry.Price, true
	} else if p.Type == PositionTypeShort {
		return p.Entry.Price - currentPrice, true
	}

	return 0.0, false
}
