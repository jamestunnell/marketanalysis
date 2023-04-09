package strategy

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/models/bar"
)

type Simple struct {
	algo   models.Algorithm
	posSig chan models.PositionType
}

const TypeSimple = "simple"

func NewSimple(algo models.Algorithm) models.Strategy {
	s := &Simple{
		algo:   algo,
		posSig: make(chan models.PositionType),
	}

	go s.waitForDirection()

	return s
}

func (s *Simple) Type() string {
	return TypeSimple
}

func (s *Simple) Params() map[string]any {
	return map[string]any{}
}

func (s *Simple) Algos() []models.Algorithm {
	return []models.Algorithm{s.algo}
}

func (s *Simple) PositionSignal() chan models.PositionType {
	return s.posSig
}

func (s *Simple) WarmupPeriod() int {
	return s.algo.WarmupPeriod()
}

func (s *Simple) WarmUp(bars []*bar.Bar) error {
	if err := s.algo.WarmUp(bars); err != nil {
		return fmt.Errorf("failed to warm up algo: %w", err)
	}

	return nil
}

func (s *Simple) Update(bar *bar.Bar) {
	s.algo.Update(bar)
}

func (s *Simple) waitForDirection() {
	for {
		dir := <-s.algo.DirectionSignal()
		switch {
		case dir > 0:
			s.posSig <- models.PositionLong
		case dir < 0:
			s.posSig <- models.PositionShort
		}
	}
}
