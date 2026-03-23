package main

func getValOrDefault[T any](val any, def ...T) (r T) {
	v, ok := val.(T)
	if !ok {
		if def != nil {
			r = def[0]
		}
		return r
	}
	return v
}

func getVal[T any](val *T) (r T) {
	if val == nil {
		return r
	}
	return *val
}

func getElm[T any](val []T, index int) (r T) {
	if index >= len(val) {
		return r
	}
	return val[index]
}
