package models

import "time"

type PositionType int

const (
	PositionLong  PositionType = 1
	PositionShort PositionType = -1
)

type Position interface {
	Entry() *TimePrice
	Exit() *TimePrice

	IsOpen() bool
	Close(t time.Time, price float64)

	OpenProfitLoss(currentPrice float64) (float64, bool)
	ClosedProfitLoss() (float64, bool)
}

type LongPosition struct {
	*PositionBase
}

type ShortPosition struct {
	*PositionBase
}

type PositionBase struct {
	entry, exit *TimePrice
}

func OpenPositions(ps []Position) []Position {
	openPs := []Position{}

	for _, p := range ps {
		if p.IsOpen() {
			openPs = append(openPs, p)
		}
	}

	return openPs
}

func ClosedPositions(ps []Position) []Position {
	closedPs := []Position{}

	for _, p := range ps {
		if !p.IsOpen() {
			closedPs = append(closedPs, p)
		}
	}

	return closedPs
}

func NewLongPosition(t time.Time, price float64) Position {
	entry := &TimePrice{
		Time:  t,
		Price: price,
	}

	return &LongPosition{
		PositionBase: &PositionBase{entry: entry, exit: nil},
	}
}

func NewShortPosition(t time.Time, price float64) Position {
	entry := &TimePrice{
		Time:  t,
		Price: price,
	}

	return &ShortPosition{
		PositionBase: &PositionBase{entry: entry, exit: nil},
	}
}

func (p *PositionBase) Entry() *TimePrice {
	return p.entry
}

func (p *PositionBase) Exit() *TimePrice {
	return p.exit
}

func (p *PositionBase) IsOpen() bool {
	return p.exit == nil
}

func (p *PositionBase) Close(t time.Time, price float64) {
	p.exit = &TimePrice{Time: t, Price: price}
}

func (p *LongPosition) ClosedProfitLoss() (float64, bool) {
	if p.IsOpen() {
		return 0.0, false
	}

	return p.exit.Price - p.entry.Price, true
}

func (p *LongPosition) OpenProfitLoss(currentPrice float64) (float64, bool) {
	if !p.IsOpen() {
		return 0.0, false
	}

	return currentPrice - p.entry.Price, true
}

func (p *ShortPosition) ClosedProfitLoss() (float64, bool) {
	if p.IsOpen() {
		return 0.0, false
	}

	return p.entry.Price - p.exit.Price, true
}

func (p *ShortPosition) OpenProfitLoss(currentPrice float64) (float64, bool) {
	if !p.IsOpen() {
		return 0.0, false
	}

	return p.entry.Price - currentPrice, true
}
