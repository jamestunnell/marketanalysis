package blocks

import (
	"reflect"
	"time"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type TypedOutputAsync[T any] struct {
	Type  string
	set   bool
	time  time.Time
	value T
	Ins   []*TypedInputAsync[T]
}

func NewTypedOutputAsync[T any]() *TypedOutputAsync[T] {
	var t T

	return &TypedOutputAsync[T]{
		Type:  reflect.TypeOf(t).String(),
		set:   false,
		value: t,
		Ins:   []*TypedInputAsync[T]{},
	}
}

func (out *TypedOutputAsync[T]) GetType() string {
	return out.Type
}

func (out *TypedOutputAsync[T]) GetConnected() []Input {
	return sliceutils.Map(out.Ins, func(in *TypedInputAsync[T]) Input {
		return in
	})
}

func (out *TypedOutputAsync[T]) IsAsynchronous() bool {
	return true
}

func (out *TypedOutputAsync[T]) IsConnected() bool {
	return len(out.Ins) != 0
}

func (out *TypedOutputAsync[T]) ClearValue() {
	out.set = false
}

func (out *TypedOutputAsync[T]) IsValueSet() bool {
	return out.set
}

func (out *TypedOutputAsync[T]) SetTimeValue(t time.Time, val T) {
	out.set = true
	out.time = t
	out.value = val
}

func (out *TypedOutputAsync[T]) SetIfConnected(t time.Time, calcVal func() T) {
	if len(out.Ins) == 0 {
		return
	}

	out.set = true
	out.time = t
	out.value = calcVal()
}

func (out *TypedOutputAsync[T]) Connect(i Input) error {
	in, ok := i.(*TypedInputAsync[T])
	if !ok {
		return commonerrs.NewErrWrongType(i.GetType(), out.Type)
	}

	out.Ins = append(out.Ins, in)
	in.out = out

	return nil
}

func (out *TypedOutputAsync[T]) DisconnectAll() {
	out.Ins = []*TypedInputAsync[T]{}
}

func (out *TypedOutputAsync[T]) GetValue() T {
	return out.value
}
