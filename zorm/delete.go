package zorm

import (
	"errors"
	"github.com/ZerQAQ/zorm/global"
	"reflect"
	"strings"
)

func (q *Operation) Delete (ptr interface{}) (int64, error) {
	val := reflect.ValueOf(ptr)
	val = global.UnpackPtr(val)

	if val.Type().Kind() != reflect.Struct {panic(errors.New("zorm: the parameter send into Delete must be struct"))}

	q.Sync(val)

	sql := "delete from " + q.table.Name
	if len(q.sqls) > 0 {
		sql += " where " + strings.Join(q.sqls, " and ")
	}

	res, err := q.driver.Database.Exec(sql, q.args...)
	if err == nil{
		id, err2 := res.RowsAffected()
		if err2 == nil {
			return id, err
		}
		return -1, err
	}

	return -1, err
}