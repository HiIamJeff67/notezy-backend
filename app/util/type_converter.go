package util

import "reflect"

func GetTypeName[T any]() string {
    var t T
	return reflect.TypeOf(t).Name()
}