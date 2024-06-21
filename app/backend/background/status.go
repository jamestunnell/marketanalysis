package background

import (
	"encoding/json"
	"time"
)

type Status struct {
	State         JobState  `json:"state"`
	Started       time.Time `json:"started"`
	Completed     time.Time `json:"completed"`
	OuterProgress float64   `json:"outerProgress"`
	InnerProgress float64   `json:"innerProgress"`
	Result        any       `json:"result"`
	ErrMsg        string    `json:"errMsg"`
}

type JobState int

const (
	Running JobState = iota
	Succeeded
	Failed
)

func (s JobState) String() string {
	switch s {
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
