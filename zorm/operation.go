package zorm

import (
	"errors"
	"github.com/ZerQAQ/zorm/global"
	"github.com/ZerQAQ/zorm/set"
	"github.com/ZerQAQ/zorm/table"
	"reflect"
	"strings"
)

type Operation struct {
	sqls   []string
	args   []interface{}
	table  *table.Table
	driver *Driver
	cols *set.Set

	offset int64
	limit int64
}

func (q *Operation) Init (driver *Driver)  {
	q.driver = driver
	q.offset = -1
	q.limit = -1

	q.sqls = make([]string, 0)
	q.args = make([]interface{}, 0)
	q.cols = set.MakeSet()
}

func (q *Operation) Sync (ptr reflect.Value) {
	typeInfo := ptr.Type()
	if q.table != nil && q.table.Name == typeInfo.Name() {
		return}
	tb, ok := q.driver.SyncTable[typeInfo.Name()]
	if !ok {panic(errors.New("zorm: table relation is not sync"))}
	q.table = &tb
}

func (q *Operation) parseArgs() {
	sqlByte := []byte(strings.Join(q.sqls, " and "))
	newSqlByte := make([]byte, 0)
	newArgs := make([]interface{}, 0)

	p := 0
	for _, elm := range sqlByte{
		if elm != '?'{
			newSqlByte = append(newSqlByte, elm)
			continue
		}

		if reflect.TypeOf(q.args[p]).Kind() != reflect.Slice {
			newSqlByte = append(newSqlByte, elm)
			newArgs = append(newArgs, q.args[p])
			p += 1
			continue
		}

		sli := reflect.ValueOf(q.args[p])
		sliceLen := sli.Len()
		for i := 0; i < sliceLen; i++ {
			if i != 0 {newSqlByte = append(newSqlByte, ',')}
			newSqlByte = append(newSqlByte, '?')

			elm := sli.Index(i)

			if global.KindIsInt(elm.Kind()) {
				newArgs = append(newArgs, elm.Int())
			} else{
				newArgs = append(newArgs, elm.String())
			}
		}
		p += 1
	}

	q.sqls = make([]string, 1)
	q.sqls[0] = string(newSqlByte)
	q.args = newArgs
}