package zorm

import (
	"orm/table"
)

type Query struct {
	sqls   []string
	args   []interface{}
	table  *table.Table
	driver *Driver

	offset int64
	limit int64
}

func (q *Query) Init (driver *Driver)  {
	q.driver = driver
	q.offset = -1
	q.limit = -1

	q.sqls = make([]string, 0)
	q.args = make([]interface{}, 0)
}