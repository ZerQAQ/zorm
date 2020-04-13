package table

import (
	"github.com/ZerQAQ/zorm/global"
	"github.com/ZerQAQ/zorm/set"
	"reflect"
	"strings"
)

func lower(v uint8) uint8 {
	if 'A' <= v && v <= 'Z' { return v - 'A' + 'a'}
	return v
}

func isLower(v uint8) bool {
	if 'a' <= v && v <= 'z' {return true} else {return false}
}

func parseName (s string) string{
	// ObjectName -> object_name
	b := []byte(s)
	b[0] = lower(b[0])
	r := make([]byte, 0)
	for i := range b {
		if isLower(b[i]) {
			r = append(r, b[i])
		} else {
			r = append(r, '_')
			r = append(r, lower(b[i]))
		}
	}
	return string(r)
}

func parseTag(s string) set.Set{
	v := strings.Split(s, " ")
	ret := set.MakeSet()
	for _, i := range v {
		ret.Insert(i)
	}
	return *ret
}

func parseInfoToRow(info reflect.StructField) Row {
	row := Row{Name:parseName(info.Name), NameRaw:info.Name,
		Pk:false, AutoIncrement:false, Null:true,
		Default:""}

	if info.Type == global.TypeInt ||
		info.Type == global.TypeInt64 {
		row.Type = "bigint"
	} else if info.Type == global.TypeInt32 {
		row.Type = "int"
	} else if info.Type == global.TypeString {
		row.Type = "varchar(256)"
	} else {
		row.Type = "undefined"
	}

	tag := parseTag(info.Tag.Get("zorm"))

	if tp, ok := tag.Like("varchar"); ok{
		row.Type = tp
	}
	if tag.Find("pk") {
		row.Pk = true
		row.Null = false
	}
	if tag.Find("auto_increment") {row.AutoIncrement = true}
	if tag.Find("not_null") {row.Null = false}

	return row
}

func parseInfoToIndex(info reflect.StructField,
	colName string) Index {
	tag := parseTag(info.Tag.Get("zorm"))

	index := Index{Unique:false,
		ColName:colName}
	if tag.Find("unique") || tag.Find("pk"){
		index.Unique = true
		index.Name = "unique_" + index.ColName
		if tag.Find("pk") {
			index.Name = "pk_" + index.ColName
		}
	} else {
		index.Name = "index_" + index.ColName
	}
	return index
}

func ParseStruct (t reflect.Type) Table {

	tableName := t.Name()

	ret := Table{Name:tableName}
	ret.Init()

	for i := 0; i < t.NumField(); i++ {
		info := t.Field(i)
		tag := parseTag(info.Tag.Get("zorm"))
		rowName := parseName(info.Name)

		if tag.Find("ignore") {continue}

		ret.Rows[rowName] = parseInfoToRow(info)
		ret.RowsName = append(ret.RowsName, rowName)

		if tag.Find("pk") || tag.Find("unique") ||
			tag.Find("index") {
			ind := parseInfoToIndex(info, rowName)
			ret.Indexs[ind.Name] = ind
		}
	}

	return ret
}