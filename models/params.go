package models

type Params map[string]Param

type Param interface {
	Type() string
	Value() any
}

func (params Params) GetInt(name string) (int, error) {
	param, found := params[name]
	if !found {
		return 0, &ErrParamNotFound{Name: name}
	}

	val, ok := param.Value().(int)
	if !ok {
		return 0, &ErrParamWrongType{Name: name, Value: param.Value()}
	}

	return val, nil
}

func (params Params) GetFloat(name string) (float64, error) {
	param, found := params[name]
	if !found {
		return 0, &ErrParamNotFound{Name: name}
	}

	val, ok := param.Value().(float64)
	if !ok {
		return 0, &ErrParamWrongType{Name: name, Value: param.Value()}
	}

	return val, nil
}
