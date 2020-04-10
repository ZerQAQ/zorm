package zorm

import (
	"database/sql"
	"errors"
	"fmt"
	"orm/global"
	"reflect"
	"strconv"
	"strings"
)

func (q *Query) parseRowToStruct (row *sql.Rows, ptr reflect.Value) {

	rowNum := len(q.table.RowsName)
	args := make([]interface{}, rowNum)
	argsPtr := make([]interface{}, rowNum)
	for i := range args {
		argsPtr[i] = &args[i]
	}

	err := row.Scan(argsPtr...)
	if err != nil {panic(err)}

	obj := ptr

	for i, rowName := range q.table.RowsName {
		rowInfo := q.table.Rows[rowName]
		field := obj.FieldByName(rowInfo.NameRaw)
		if field.IsValid() {
			tp := field.Type()
			if tp == global.TypeInt32 || tp == global.TypeInt64 ||
				tp == global.TypeInt {
				field.SetInt(global.ParseInt64(args[i]))
			} else {
				field.SetString(string(args[i].([]byte)))
			}
		}
	}
}

func (q *Query) Sync (ptr reflect.Value) {
	typeInfo := ptr.Type()
	tb, ok := q.driver.SyncTable[typeInfo.Name()]
	if !ok {panic(errors.New("zorm: table relation is not sync"))}
	q.table = &tb
}

func (q *Query) Get (ptr interface{}) bool {

	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {panic("zorm: parameter send to Get must be pointer")}

	q.Sync(reflect.ValueOf(ptr).Elem())

	sql := "select * from " + q.table.Name
	if len(q.sqls) > 0 {sql += " where " + strings.Join(q.sqls, " and ")}
	sql += " limit 1"

	//fmt.Println(sql, q.args)

	res, err := q.driver.Database.Query(sql, q.args...)

	if err != nil {panic(err)}

	if !res.Next() {return false}
	q.parseRowToStruct(res, reflect.ValueOf(ptr).Elem())
	return true
}

/*
	Find返回多个数据，获得的数据个数不会超过slice的长度
 */
func (q *Query) Find (sli interface{}) error {
	if reflect.TypeOf(sli).Kind() != reflect.Ptr {panic("zorm: the parameter put into Find must be pointer")}
	val := reflect.ValueOf(sli).Elem()
	if val.Kind() != reflect.Slice {panic("zorm: the parameter put into Find must be slice pointer")}
	//if int64(val.Len()) < q.limit {return errors.New("zorm: the len of slice must be greater than limit")}
	sliceLen := int64(val.Len())
	if sliceLen == 0 {return nil}
	// limit 取两者的最小值
	if q.limit == -1 || sliceLen < q.limit {q.limit = sliceLen}

	firval := val.Index(0)
	if firval.Type().Kind() == reflect.Ptr {firval = firval.Elem()}
	q.Sync(firval)

	sql := "select * from " + q.table.Name
	if len(q.sqls) > 0 {sql += " where " + strings.Join(q.sqls, " and ")}
	if q.limit > 0 {sql += " limit " + strconv.Itoa(int(q.limit))}
	if q.offset >= 0 {sql += " offset " + strconv.Itoa(int(q.offset))}

	//fmt.Println(sql, q.args)

	res, err := q.driver.Database.Query(sql, q.args...)
	/*
		sql返回错误有可能是sql语句写错了，应该马上被得到改正，所以我用了panic
	 */
	if err != nil {panic(err)}

	var p int64
	p = 0
	for res.Next() {
		valElm := val.Index(int(p)); p++
		if valElm.Type().Kind() == reflect.Ptr {valElm = valElm.Elem()}

		q.parseRowToStruct(res, valElm)
		if p >= sliceLen {return nil}
	}

	return nil
}

func (q *Query) Count (ptr interface{}) int {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {panic("zorm: parameter send to Count must be pointer")}

	q.Sync(reflect.ValueOf(ptr).Elem())

	sql := "select * from " + q.table.Name
	if len(q.sqls) > 0 {sql += " where " + strings.Join(q.sqls, " and ")}
	if q.limit > 0 {sql += " limit " + strconv.Itoa(int(q.limit))}
	if q.offset >= 0 {sql += " offset " + strconv.Itoa(int(q.offset))}
	sql = "select count(*) from (" + sql + ") as t"
	fmt.Println(sql)

	res, err := q.driver.Database.Query(sql, q.args...)
	if err != nil {panic(err)}

	res.Next()
	var num int
	res.Scan(&num)

	return num
}