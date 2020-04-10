package global

import (
	"reflect"
)

var TypeInt64 = reflect.TypeOf(*new(int64))
var TypeInt32 = reflect.TypeOf(*new(int32))
var TypeInt = reflect.TypeOf(*new(int))

var TypeString = reflect.TypeOf(*new(string))

var TypePtr = reflect.Ptr