package zorm

import (
	"fmt"
	"log"
	"orm/table"
)
func (d *Driver) initElm(t *table.Table, name string,
	rown int64, cmd string,
	f func(*table.Table, []string))  {

	rows, err := d.db.Query(cmd)
	if err != nil {panic(err)}

	for rows.Next() {
		//构造ifs传入scan
		var val = make([][]byte, rown)
		var ifs = make([]interface{}, rown)
		for i := range val {
			ifs[i] = &val[i]
		}

		err := rows.Scan(ifs...)
		if err != nil {log.Fatal(err)}

		var strval = make([]string, rown)
		for i, elm := range val {
			strval[i] = string(elm)
			if rown == 6 {
				fmt.Println(i, ":", strval[i])
			}
		}
		f(t, strval)
	}
}

func (d *Driver) initRows(t *table.Table, name string) {
	d.initElm(t, name, 6, "desc " + name,
		func(t *table.Table, s []string) {
			t.Rows[s[0]] = table.MakeRow(s)
		})
}

func (d *Driver) initIndex(t *table.Table, name string){
	d.initElm(t, name, 15, "show index from " + name,
		func(t *table.Table, s []string) {
			ind := table.MakeIndex(s)
			t.Indexs[ind.Name] = ind
		})
}

//将数据库中名字为name的表转化成Table对象
func (d *Driver) makeTable (name string) table.Table {
	var ret = table.Table{Name:name}
	ret.Init()
	d.initRows(&ret, name)
	d.initIndex(&ret, name)
	return ret
}