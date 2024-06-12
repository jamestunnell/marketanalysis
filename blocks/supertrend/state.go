package supertrend

type State struct {
	parent *Supertrend
}

func NewState(parent *Supertrend) *State {
	return &State{parent: parent}
}

// func (s *State) isInputSet() bool {
// 	return s.parent.in.IsValueSet()
// }

// func (s *State) isInputAboveOpenThresh() bool {
// 	return s.parent.in.GetValue() > s.parent.openThresh.CurrentVal
// }

// func (s *State) isInputBelowNegOpenThresh() bool {
// 	return s.parent.in.GetValue() < -s.parent.openThresh.CurrentVal
// }

// func (s *State) isInputBelowCloseThresh() bool {
// 	return s.parent.in.GetValue() < s.parent.closeThresh.CurrentVal
// }

// func (s *State) isInputAboveNegCloseThresh() bool {
// 	return s.parent.in.GetValue() > -s.parent.closeThresh.CurrentVal
// }

// func (s *State) setOutput(val float64) {
// 	s.parent.out.SetValue(val)
// }
