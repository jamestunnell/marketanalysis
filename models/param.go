package models

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type Params map[string]Param

type Param interface {
	Type() string
	Constraint() Constraint
	LoadVal([]byte) error
	StoreVal() ([]byte, error)
	GetVal() any
	SetVal(any) error
}

type TypedParam[T any] struct {
	Value  T
	Constr Constraint
}

func NewParam[T any](c Constraint) *TypedParam[T] {
	var t T

	return &TypedParam[T]{
		Value:  t,
		Constr: c,
	}
}

func (p *TypedParam[T]) Type() string {
	return reflect.TypeOf(p.Value).String()
}

func (p *TypedParam[T]) Constraint() Constraint {
	return p.Constr
}

func (p *TypedParam[T]) LoadVal(d []byte) error {
	var val T

	if err := json.Unmarshal(d, &val); err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	p.Value = val

	return nil
}

func (p *TypedParam[T]) StoreVal() ([]byte, error) {
	d, err := json.Marshal(p.Value)
	if err != nil {
		return []byte{}, fmt.Errorf("marshal failed: %w", err)
	}

	return d, nil
}

func (p *TypedParam[T]) GetVal() any {
	return p.Value
}

func (p *TypedParam[T]) SetVal(val any) error {
	if p.Constr != nil {
		if err := p.Constr.Check(val); err != nil {
			return fmt.Errorf("constraint failed on value %v: %w", val, err)
		}
	}

	tVal, ok := val.(T)
	if !ok {
		actual := reflect.TypeOf(val).String()

		return commonerrs.NewErrWrongType(actual, p.Type())
	}

	p.Value = tVal

	return nil
}
