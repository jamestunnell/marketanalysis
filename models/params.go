package models

type Params map[string]any

func (params Params) GetInt(name string) (int, error) {
	obj, found := params[name]
	if !found {
		return 0, &ErrParamNotFound{Name: name}
	}

	val, ok := obj.(int)
	if !ok {
		return 0, &ErrParamWrongType{Name: name}
	}

	return val, nil
}

func (params Params) GetFloat(name string) (float64, error) {
	obj, found := params[name]
	if !found {
		return 0, &ErrParamNotFound{Name: name}
	}

	val, ok := obj.(float64)
	if !ok {
		return 0, &ErrParamWrongType{Name: name}
	}

	return val, nil
}
