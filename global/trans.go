package global

import "reflect"

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