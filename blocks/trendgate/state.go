package trendgate

type State struct {
	parent *TrendGate
}

func NewState(parent *TrendGate) *State {
	return &State{parent: parent}
}

func (s *State) OpenThresh() float64 {
	return s.parent.openThresh.CurrentVal
}

func (s *State) CloseThresh() float64 {
	return s.parent.closeThresh.CurrentVal
}

func (s *State) DebouncePeriod() int {
	return s.parent.debouncePeriod.CurrentVal
}
