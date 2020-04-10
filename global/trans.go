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