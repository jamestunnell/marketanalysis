package models

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type Input interface {
	GetType() string
	Connect(Output) error
	IsConnected() bool
}

type Inputs map[string]Input

func NewTypedInput[T any]() *TypedInput[T] {
	var t T

	return &TypedInput[T]{
		Type: reflect.TypeOf(t).String(),
		Out: nil,
	}
}

type TypedInput[T any] struct {
	Type string
	Out *TypedOutput[T]
}

func (in *TypedInput[T]) GetType() string {
	return in.Type
}

func (in *TypedInput[T]) Connect(o Output) error {
	out, ok := o.(*TypedOutput[T])
	if !ok {
		return commonerrs.NewErrWrongType(o.GetType(), in.Type)
	}

	in.Out = out

	return nil
}

func (in *TypedInput[T]) IsConnected() bool {
	return in.Out != nil
}
