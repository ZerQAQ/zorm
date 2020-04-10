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

func (q *Query) Where (cmd string, args ...interface{}) *Query {
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

func (q *Query) Id (id int64) *Query {
	return q.Where("id = ?", id)
}

func (q *Query) Limit (offset int64, limit int64) *Query {
	q.offset = offset
	q.limit = limit
	return q
}