package zorm

import (
	"database/sql"
	"errors"
	"github.com/ZerQAQ/zorm/global"
	"reflect"
	"strconv"
	"strings"
)

func (q *Operation) parseRowToStruct (row *sql.Rows, ptr reflect.Value) {

	rowNum := len(q.table.RowsName)
	args := make([]interface{}, rowNum)
	argsPtr := make([]interface{}, rowNum)
	for i := range args {
		argsPtr[i] = &args[i]
	}

	err := row.Scan(argsPtr...)
	if err != nil {panic(err)}

	//fmt.Println(args)

	obj := ptr

	for i, rowName := range q.table.RowsName {
		rowInfo := q.table.Rows[rowName]
		field := obj.FieldByName(rowInfo.NameRaw)
		if field.IsValid() {
			if global.KindIsInt(field.Kind()){
				// 两种情况，用Find的时候int会被Scan成[]byte
				if global.KindIsInt(reflect.TypeOf(args[i]).Kind()) {
					field.SetInt(global.ParseInt64(args[i]))
				} else {
					v, _ := strconv.Atoi(string(args[i].([]byte)))
					field.SetInt(int64(v))
				}
			} else {
				field.SetString(string(args[i].([]byte)))
			}
		}
	}
}

func (q *Operation) Get (ptr interface{}) (bool, error) {

	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {panic(errors.New("zorm: parameter send to Get must be pointer"))}
	val := global.UnpackPtr(reflect.ValueOf(ptr))

	if val.Type().Kind() != reflect.Struct {panic(errors.New("zorm: parameter send to Get must be a struct pointer"))}

	q.Sync(val)

	sql := "select * from " + q.table.Name
	if len(q.sqls) > 0 {sql += " where " + strings.Join(q.sqls, " and ")}
	sql += " limit 1"

	//fmt.Println(sql, q.args)

	res, err := q.driver.Database.Query(sql, q.args...)

	if err != nil {return false, err}

	if !res.Next() {return false, nil}
	q.parseRowToStruct(res, val)
	return true, nil
}

/*
	Find返回多个数据，获得的数据个数不会超过slice的长度
 */
func (q *Operation) Find (sli interface{}) (int64, error) {
	if reflect.TypeOf(sli).Kind() != reflect.Ptr {panic("zorm: the parameter put into Find must be pointer")}
	val := reflect.ValueOf(sli).Elem()
	if val.Kind() != reflect.Slice {panic("zorm: the parameter put into Find must be slice pointer")}
	//if int64(val.Len()) < q.limit {return errors.New("zorm: the len of slice must be greater than limit")}
	sliceLen := int64(val.Len())
	if sliceLen == 0 {return 0, nil}
	// limit 取两者的最小值
	if q.limit == -1 || sliceLen < q.limit {q.limit = sliceLen}

	firVal := val.Index(0)
	if firVal.Type().Kind() == reflect.Ptr {
		firVal = firVal.Elem()}
	q.Sync(firVal)

	sql := "select * from " + q.table.Name
	if len(q.sqls) > 0 {sql += " where " + strings.Join(q.sqls, " and ")}
	if q.limit > 0 {sql += " limit " + strconv.Itoa(int(q.limit))}
	if q.offset >= 0 {sql += " offset " + strconv.Itoa(int(q.offset))}

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
		if p >= sliceLen {return sliceLen, nil}
	}

	return p, nil
}

func (q *Operation) Count (ptr interface{}) (int64, error) {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {panic("zorm: parameter send to Count must be pointer")}

	q.Sync(reflect.ValueOf(ptr).Elem())

	sql := "select * from " + q.table.Name
	if len(q.sqls) > 0 {sql += " where " + strings.Join(q.sqls, " and ")}
	if q.limit > 0 {sql += " limit " + strconv.Itoa(int(q.limit))}
	if q.offset >= 0 {sql += " offset " + strconv.Itoa(int(q.offset))}
	sql = "select count(*) from (" + sql + ") as t"

	res, err := q.driver.Database.Query(sql, q.args...)

	res.Next()
	var num int
	res.Scan(&num)

	return int64(num), err
}