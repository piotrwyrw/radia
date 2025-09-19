package util

import "reflect"

func IsInteger(t reflect.Type) bool {
	return t.Kind() >= reflect.Int && t.Kind() <= reflect.Int64
}

func IsFloat(t reflect.Type) bool {
	return t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64
}
