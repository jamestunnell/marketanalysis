package blocks

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type Input interface {
	GetType() string

	IsAsynchronous() bool
	IsConnected() bool
	IsOptional() bool
	IsValueSet() bool

	Connect(Output) error
}

type TypedInput[T any] struct {
	Type     string
	Optional bool

	out *TypedOutput[T]
}

type Inputs map[string]Input

func NewTypedInput[T any]() *TypedInput[T] {
	var val T

	return &TypedInput[T]{
		Type:     reflect.TypeOf(val).String(),
		Optional: false,
		out:      nil,
	}
}

func NewTypedInputOptional[T any]() *TypedInput[T] {
	var val T

	return &TypedInput[T]{
		Type:     reflect.TypeOf(val).String(),
		Optional: true,
		out:      nil,
	}
}

func (in *TypedInput[T]) GetType() string {
	return in.Type
}

func (in *TypedInput[T]) IsAsynchronous() bool {
	return false
}

func (in *TypedInput[T]) IsConnected() bool {
	return in.out != nil
}

func (in *TypedInput[T]) IsOptional() bool {
	return in.Optional
}

func (in *TypedInput[T]) IsValueSet() bool {
	return in.IsConnected() && in.out.IsValueSet()
}

func (in *TypedInput[T]) Connect(o Output) error {
	out, ok := o.(*TypedOutput[T])
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

func (in *TypedInput[T]) GetValue() T {
	return in.out.value
}
