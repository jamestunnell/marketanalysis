package models

import "reflect"

type Output interface {
	GetType() string
	GetValue() any
}

type Outputs map[string]Output

type TypedOutput[T any] struct {
	Value T
	Type  string
}

func NewTypedOutput[T any]() *TypedOutput[T] {
	var t T

	return &TypedOutput[T]{
		Type:  reflect.TypeOf(t).String(),
		Value: t,
	}
}

func (out *TypedOutput[T]) GetType() string {
	return out.Type
}

func (out *TypedOutput[T]) GetValue() any {
	return out.Value
}
