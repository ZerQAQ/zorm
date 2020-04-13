package global

import (
	"reflect"
	"unsafe"
)

func ParseInt64 (v interface{}) int64 {
	tp := reflect.TypeOf(v)
	if tp == TypeInt64 {
		return v.(int64)
	} else if tp == TypeInt32 {
		return int64(v.(int32))
	} else if tp == TypeInt {
		return int64(v.(int))
	} else {
		return 0
	}
}

func UnpackPtr(value reflect.Value) reflect.Value {
	if value.Type().Kind() == reflect.Ptr {return value.Elem()} else{return value}
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}