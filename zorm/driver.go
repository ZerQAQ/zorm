package zorm

import (
	"database/sql"
	"orm/set"
	"orm/table"
)

type Driver struct {
	Database  *sql.DB
	tableSet  *set.Set
	SyncTable map[string]table.Table
}

func Open(name string, sour string) *Driver {
	ret := new(Driver)
	ret.Connect(name, sour)
	ret.init()
	return ret
}

func (d *Driver) Sync (s interface{}) bool {
	// in是用户输入的表结构，tb是数据库中已经存在的表结构
	in := table.ParseStruct(s)
	if d.tableSet.Find(in.Name) {
		tb := d.makeTable(in.Name)

		if !table.IsContain(&tb, &in){
			panic("Conflict, cant sync.")
		}

		if table.IsContain(&in, &tb){
			d.SyncTable[in.Name] = in
			return true
		}

		// in 和 tb 的差距
		ad := table.Sub(&in, &tb)
		//fmt.Println("altering")
		in.RowsName = tb.RowsName

		in.UpdateRowList(&ad)
		d.alterTable(&ad)
		//debug.PrtTable(&ad)
	} else {
		//fmt.Println("Creating")
		d.createTable(&in)
		//debug.Prt(&in)
	}
	d.SyncTable[in.Name] = in
	return true
}

func (d *Driver) Count (ptr interface{}) int {
	q := new(Query)
	q.Init(d)
	return q.Count(ptr)
}

func (d *Driver) Where (cmd string, args ...interface{}) *Query {
	q := new(Query)
	q.Init(d)
	return q.Where(cmd, args...)
}

func (d *Driver) Id (id int64) *Query {
	q := new(Query)
	q.Init(d)
	return q.Id(id)
}