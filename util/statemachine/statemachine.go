package statemachine

import (
	"time"

	"github.com/rs/zerolog/log"
)

type StateMachine[T any] struct {
	Name    string
	Current State[T]
}

func New[T any](name string, start State[T]) *StateMachine[T] {
	log.Debug().
		Str("startState", start.GetName()).
		Msgf("%s state machine: starting", name)

	start.Enter()

	return &StateMachine[T]{
		Current: start,
		Name:    name,
	}
}

func (sm *StateMachine[T]) Run(t time.Time, val T) {
	if next := sm.Current.Run(t, val); next != nil {
		log.Trace().
			Str("nextState", next.GetName()).
			Stringer("timestamp", t).
			Msgf("%s state machine: changing state", sm.Name)

		sm.Current.Exit()

		next.Enter()

		sm.Current = next
	}
}
