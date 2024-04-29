package blocks

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type Output interface {
	GetType() string
	GetConnected() []Input

	IsConnected() bool
	IsValueSet() bool

	Connect(Input) error
	DisconnectAll()
}

type Outputs map[string]Output

type TypedOutput[T any] struct {
	Type  string
	set   bool
	value T
	Ins   []*TypedInput[T]
}

func NewTypedOutput[T any]() *TypedOutput[T] {
	var t T

	return &TypedOutput[T]{
		Type:  reflect.TypeOf(t).String(),
		set:   false,
		value: t,
		Ins:   []*TypedInput[T]{},
	}
}

func (out *TypedOutput[T]) GetType() string {
	return out.Type
}

func (out *TypedOutput[T]) GetConnected() []Input {
	return sliceutils.Map(out.Ins, func(in *TypedInput[T]) Input {
		return in
	})
}

func (out *TypedOutput[T]) IsConnected() bool {
	return len(out.Ins) != 0
}

func (out *TypedOutput[T]) IsValueSet() bool {
	return out.set
}

func (out *TypedOutput[T]) SetValue(val T) {
	out.set = true
	out.value = val
}

func (out *TypedOutput[T]) SetIfConnected(calcVal func() T) {
	if len(out.Ins) == 0 {
		return
	}

	out.set = true
	out.value = calcVal()
}

func (out *TypedOutput[T]) Connect(i Input) error {
	in, ok := i.(*TypedInput[T])
	if !ok {
		return commonerrs.NewErrWrongType(i.GetType(), out.Type)
	}

	out.Ins = append(out.Ins, in)
	in.out = out

	return nil
}

func (out *TypedOutput[T]) DisconnectAll() {
	out.Ins = []*TypedInput[T]{}
}

func (out *TypedOutput[T]) GetValue() T {
	return out.value
}
