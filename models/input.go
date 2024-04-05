package models

import (
	"reflect"
)

type Input interface {
	GetType() string

	IsConnected() bool
}

type Inputs map[string]Input

type TypedInput[T any] struct {
	Type string

	out            *TypedOutput[T]
	value          T
	connected, set bool
}

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

func (in *TypedInput[T]) IsSet() bool {
	return in.set
}

func (in *TypedInput[T]) Connect() {
	in.connected = true
}

func (in *TypedInput[T]) Disconnect() {
	in.connected = false
}

func (in *TypedInput[T]) Set(val T) {
	in.value = val
	in.set = true
}

func (in *TypedInput[T]) Get() T {
	return in.value
}

func (in *TypedInput[T]) Reset() {
	var val T

	in.value = val
	in.set = false
}
