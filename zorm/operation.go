package zorm

import (
	"errors"
	"orm/table"
	"reflect"
)

type Operation struct {
	sqls   []string
	args   []interface{}
	table  *table.Table
	driver *Driver

	offset int64
	limit int64
}

func (q *Operation) Init (driver *Driver)  {
	q.driver = driver
	q.offset = -1
	q.limit = -1

	q.sqls = make([]string, 0)
	q.args = make([]interface{}, 0)
}

func (q *Operation) Sync (ptr reflect.Value) {
	typeInfo := ptr.Type()
	if q.table != nil && q.table.Name == typeInfo.Name() {
		return}
	tb, ok := q.driver.SyncTable[typeInfo.Name()]
	if !ok {panic(errors.New("zorm: table relation is not sync"))}
	q.table = &tb
}