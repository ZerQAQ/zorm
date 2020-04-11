package zorm

import (
	"log"
	"github.com/ZerQAQ/zorm/table"
)

// 巧妙的做法
func (d *Driver) initElm(t *table.Table, name string,
	rowNum int64, cmd string,
	f func(*table.Table, []string))  {

	rows, err := d.Database.Query(cmd)
	if err != nil {panic(err)}

	for rows.Next() {
		//构造ifs传入scan
		var val = make([][]byte, rowNum)
		var ifs = make([]interface{}, rowNum)
		for i := range val {
			ifs[i] = &val[i]
		}

		err := rows.Scan(ifs...)
		if err != nil {log.Fatal(err)}

		var strVal = make([]string, rowNum)
		for i, elm := range val {
			strVal[i] = string(elm)
		}
		f(t, strVal)
	}
}

func (d *Driver) initRows(t *table.Table, name string) {
	d.initElm(t, name, 6, "desc " + name,
		func(t *table.Table, s []string) {
			row := table.MakeRow(s)
			t.Rows[row.Name] = row
			t.RowsName = append(t.RowsName, row.Name)
		})
}

func (d *Driver) initIndex(t *table.Table, name string){
	d.initElm(t, name, 15, "show index from " + name,
		func(t *table.Table, s []string) {
			ind := table.MakeIndex(s)
			if t.Rows[ind.ColName].Pk {
				ind.Name = "pk_" + ind.ColName
			}
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