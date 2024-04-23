package models

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type Input interface {
	GetType() string

	IsConnected() bool
	IsValueSet() bool

	Connect(Output) error
}

type TypedInput[T any] struct {
	Type string

	out   *TypedOutput[T]
	value T
}

type Inputs map[string]Input

func NewTypedInput[T any]() *TypedInput[T] {
	var val T

	return &TypedInput[T]{
		Type:  reflect.TypeOf(val).String(),
		out:   nil,
		value: val,
	}
}

func (in *TypedInput[T]) GetType() string {
	return in.Type
}

func (in *TypedInput[T]) IsConnected() bool {
	return in.out != nil
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
