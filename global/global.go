package global

import (
	"reflect"
)

var typeInt64 = reflect.TypeOf(*new(int64))
var typeInt32 = reflect.TypeOf(*new(int32))
var typeInt = reflect.TypeOf(*new(int))

var typeString = reflect.TypeOf(*new(string))