package zorm

import (
	"database/sql"
	"errors"
	"fmt"
	"orm/global"
	"reflect"
)

func (q *Operation) parseStructToArgs(val reflect.Value) []interface{} {
	args := make([]interface{}, 0)
	for _, elm := range q.table.RowsName{
		rowNameRaw := q.table.Rows[elm].NameRaw
		field := val.FieldByName(rowNameRaw)
		// struct中缺少某个表中的字段
		if !field.IsValid() {panic(errors.New("zorm: the struct is not sync"))}
		if global.KindIsInt(field.Type().Kind()) {
			args = append(args, field.Int())
		} else if field.Type().Kind() == reflect.String {
			args = append(args, field.String())
		}
	}
	return args
}

func (q *Operation) insertOne (val reflect.Value, cursor *sql.Tx) (sql.Result, error) {
	q.Sync(val)
	/*
		解析args的时候可能会panic，添加defer保证回滚
	 */
	defer cursor.Rollback()

	sql := "insert into " + q.table.Name + " value ("
	for i, _ := range q.table.RowsName {
		if i == 0 {sql += " ? "} else{
			sql += ",? "
		}
	}
	sql += ")"

	args := q.parseStructToArgs(val)
	fmt.Println(sql, args)
	return cursor.Exec(sql, args...)
}

func (q *Operation) insert (ptrs ...interface{}) (int64, error) {
	cursor, err := q.driver.Database.Begin()
	if err != nil {return -1, err}
	var res sql.Result
	for _, elm := range ptrs {
		val := reflect.ValueOf(elm)
		if val.Type().Kind() == reflect.Ptr {val = val.Elem()}
		res, err = q.insertOne(val, cursor)
		if err != nil {
			//若中途出错，回滚所有操作
			cursor.Rollback()
			return -1, err
		}
	}
	cursor.Commit()
	id, err := res.LastInsertId()
	if err == nil {return id, nil} else{
		return -1, nil
	}
}

/*
	传入一个slice或者slicePtr，插入多个数据
 */
func (q *Operation) insertMany (sli interface{}) (int64, error) {
	val := reflect.ValueOf(sli)
	val = global.UnpackPtr(val)
	if val.Type().Kind() != reflect.Slice {panic("zorm: the parameter send into InsertMany must be slice")}
	sliceLen := val.Len()
	if sliceLen == 0 {return -1, nil}

	q.Sync(global.UnpackPtr(val.Index(0)))

	sql := "insert into " + q.table.Name + " value "

	subsql := "( "
	for i := 0; i <sliceLen; i++ {
		if i == 0 {subsql += " ? "} else {subsql += ",? "}
	}
	subsql += " )\n"

	args := make([]interface{}, 0)
	for i := 0; i < sliceLen; i++  {
		elm := global.UnpackPtr(val.Index(i))
		args = append(args, q.parseStructToArgs(elm)...)
		if i != 0 {sql += ","}
		sql += subsql
	}
	res, err := q.driver.Database.Exec(sql, args...)

	if err == nil{
		id, err2 := res.LastInsertId()
		id += int64(sliceLen) - 1
		if err2 != nil {id = -1}
		return id, nil
	}

	return -1, err
}