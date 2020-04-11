package zorm

import (
	"errors"
)

func countNum (s string) int64 {
	bs := []byte(s)
	var ret int64
	ret = 0
	for _, elm := range bs{
		if elm == '?' {ret += 1}
	}
	return ret
}

func (q *Operation) Where (cmd string, args ...interface{}) *Operation {
	if cmd == "" {return q}
	if countNum(cmd) != int64(len(args)) {
		panic(errors.New("zorm: the number of '?' in query command should be equal to the args number"))
	}

	addition := "( "
	addition += cmd
	addition += " )"

	q.sqls = append(q.sqls, addition)
	q.args = append(q.args, args...)

	return q
}

func (q *Operation) Id (id int64) *Operation {
	return q.Where("id = ?", id)
}

func (q *Operation) Limit (offset int64, limit int64) *Operation {
	q.offset = offset
	q.limit = limit
	return q
}

func (q *Operation) Col (colName ...string) *Operation {
	for _, elm := range colName {
		q.cols.Insert(elm)
	}
	return q
}