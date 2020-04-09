package orm

import (
	"log"
	"orm/table"
)

func (d *Drive) initElm(t *table.Table, name string, rown int64,
	f func(*table.Table, []string))  {

	rows, err := d.db.Query("desc " + name)
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
		}
		f(t, strval)
		t.Rows[strval[0]] = table.MakeRow(strval)
	}
}

func (d *Drive) initRows(t *table.Table, name string) {
	d.initElm(t, name, 6,
		func(t *table.Table, s []string) {
			t.Rows[s[0]] = table.MakeRow(s)
		})
}

func (d *Drive) initIndex(t *table.Table, name string){
	d.initElm(t, name, 15,
		func(t *table.Table, s []string) {
			t.Indexs[s[0]] = table.MakeIndex(s)
		})
}

func (d *Drive) makeTable (name string) table.Table {
	var ret = table.Table{Name:name}
	d.initRows(&ret, name)
	d.initIndex(&ret, name)
	return ret
}