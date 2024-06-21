package background

import (
	"encoding/json"
	"time"
)

type Status struct {
	State         JobState  `json:"state"`
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`
	OuterProgress float64   `json:"outerProgress"`
	InnerProgress float64   `json:"innerProgress"`
	Result        Cloneable `json:"result"`
	ErrMsg        string    `json:"errMsg"`
}

type Cloneable interface {
	Clone() Cloneable
}

type JobState int

const (
	Queued JobState = iota
	Running
	Succeeded
	Failed
)

func (s *Status) Clone() *Status {
	return &Status{
		State:  s.State,
		Start:  s.Start,
		End:    s.End,
		Result: s.Result.Clone(),
		ErrMsg: s.ErrMsg,
	}
}

func (s JobState) String() string {
	switch s {
	case Queued:
		return "queued"
	case Running:
		return "running"
	case Succeeded:
		return "succeeded"
	case Failed:
		return "failed"
	}

	return ""
}

func (s JobState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
