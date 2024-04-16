package models

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type Output interface {
	GetType() string

	Connect(Input) error
	DisconnectAll()
}

type Outputs map[string]Output

type TypedOutput[T any] struct {
	Type     string
	Value    T
	Prev     T
	SetCount int
	Ins      []*TypedInput[T]
}

func NewTypedOutput[T any]() *TypedOutput[T] {
	var t T

	return &TypedOutput[T]{
		Type:     reflect.TypeOf(t).String(),
		Value:    t,
		Prev:     t,
		SetCount: 0,
		Ins:      []*TypedInput[T]{},
	}
}

func (out *TypedOutput[T]) GetType() string {
	return out.Type
}

func (out *TypedOutput[T]) SetIfConnected(calcVal func() T) {
	if len(out.Ins) == 0 {
		return
	}

	out.Set(calcVal())
}

func (out *TypedOutput[T]) Set(val T) {
	if out.SetCount > 0 {
		out.Prev = out.Value
	}

	for _, in := range out.Ins {
		in.Set(val)
	}

	out.Value = val
	out.SetCount++
}

func (out *TypedOutput[T]) Connect(i Input) error {
	in, ok := i.(*TypedInput[T])
	if !ok {
		return commonerrs.NewErrWrongType(i.GetType(), out.Type)
	}

	out.Ins = append(out.Ins, in)

	return nil
}

func (out *TypedOutput[T]) DisconnectAll() {
	out.Ins = []*TypedInput[T]{}
}
