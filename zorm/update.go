package zorm

import (
	"github.com/ZerQAQ/zorm/global"
	"reflect"
	"strings"
)

/*
	返回的是修改的行数以及错误
 */
func (q *Operation) Update (ptr interface{}) (int64, error) {
	val := global.UnpackPtr(reflect.ValueOf(ptr))
	q.Sync(val)

	subSql := ""

	fir := true

	args := make([]interface{}, 0)
	for _, rowName := range q.table.RowsName{
		rowNameRaw := q.table.Rows[rowName].NameRaw
		field := val.FieldByName(rowNameRaw)
		if !field.IsValid() {continue}

		if global.KindIsInt(field.Type().Kind()) {
			v := field.Int()
			//空值默认不更新 除非指定了特定的行
			if v == 0 && !q.cols.Find(rowName) {continue}
			args = append(args, v)
		} else {
			v := field.String()
			if v == "" && !q.cols.Find(rowName) {continue}
			args = append(args, v)
		}

		if !fir {
			subSql += ","}
		subSql += rowName + " = ? "
		fir = false
	}

	args = append(args, q.args...)
	sql := "update " + q.table.Name + " set " + subSql + " where " + strings.Join(q.sqls, " and ")

	res, err := q.driver.Database.Exec(sql, args...)
	if err == nil {
		id, err2 := res.RowsAffected()
		if err2 == nil { return id, err }
		return -1, err
	}
	return -1, err
}