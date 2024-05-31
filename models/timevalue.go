package models

import "time"

type TimeValue[TValue any] struct {
	Time  time.Time `json:"t"`
	Value TValue    `json:"v"`
}

func NewTimeValue[TValue any](t time.Time, v TValue) TimeValue[TValue] {
	return TimeValue[TValue]{
		Time:  t,
		Value: v,
	}
}
