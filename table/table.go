package table

import (

)

type Table struct {
	Name string
	Rows map[string]Row
	Indexs map[string]Index
}

func ParseStruct (s interface{}) Table {
	
}