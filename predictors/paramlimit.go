package predictors

type UpperLimit interface {
	Change(val any) bool
}

type TypedUpperLimit[T any] struct {
	Value T
}

func NewParamLimit[T any](val T) *TypedUpperLimit[T] {
	return &TypedUpperLimit[T]{Value: val}
}

func (lim *TypedUpperLimit[T]) Change(val any) bool {
	tVal, ok := val.(T)
	if !ok {
		return false
	}

	lim.Value = tVal

	return true
}
