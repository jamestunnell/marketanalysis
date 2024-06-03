package blocks

import (
	"reflect"
	"time"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type TypedInputAsync[T any] struct {
	Type     string
	Optional bool

	out *TypedOutputAsync[T]
}

func NewTypedInputAsync[T any]() *TypedInputAsync[T] {
	var val T

	return &TypedInputAsync[T]{
		Type:     reflect.TypeOf(val).String(),
		Optional: false,
		out:      nil,
	}
}

func NewTypedInputAsyncOptional[T any]() *TypedInputAsync[T] {
	var val T

	return &TypedInputAsync[T]{
		Type:     reflect.TypeOf(val).String(),
		Optional: true,
		out:      nil,
	}
}

func (in *TypedInputAsync[T]) GetType() string {
	return in.Type
}

func (in *TypedInputAsync[T]) IsAsynchronous() bool {
	return true
}

func (in *TypedInputAsync[T]) IsConnected() bool {
	return in.out != nil
}

func (in *TypedInputAsync[T]) IsOptional() bool {
	return in.Optional
}

func (in *TypedInputAsync[T]) IsValueSet() bool {
	return in.IsConnected() && in.out.IsValueSet()
}

func (in *TypedInputAsync[T]) Connect(o Output) error {
	out, ok := o.(*TypedOutputAsync[T])
	if !ok {
		return commonerrs.NewErrWrongType(o.GetType(), in.Type)
	}

	out.Ins = append(out.Ins, in)

	in.out = out

	return nil
}

// func (in *TypedInput[T]) Disconnect() {
// 	in.connected = false
// }

func (in *TypedInputAsync[T]) GetTime() time.Time {
	return in.out.time
}

func (in *TypedInputAsync[T]) GetValue() T {
	return in.out.value
}
