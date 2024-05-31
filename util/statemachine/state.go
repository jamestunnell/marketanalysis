package statemachine

import "time"

type State[T any] interface {
	GetName() string

	Enter()
	Run(t time.Time, val T) State[T]
	Exit()
}
