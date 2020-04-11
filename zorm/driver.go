package zorm

import (
	"database/sql"
	"github.com/ZerQAQ/zorm/set"
	"github.com/ZerQAQ/zorm/table"
)

type Driver struct {
	Database  *sql.DB
	tableSet  *set.Set
	syncTable map[string]table.Table
}

func Open(name string, sour string) (*Driver, error) {
	ret := new(Driver)
	err := ret.Connect(name, sour)
	ret.init()
	return ret, err
}

func (d *Driver) InsertMany (sli interface{}) (int64, error) {
	q := new(Operation)
	q.Init(d)
	return q.insertMany(sli)
}

func (d *Driver) Col (ptrs ...string) *Operation {
	q := new(Operation)
	q.Init(d)
	return q.Col(ptrs...)
}

func (d *Driver) Insert (ptrs ...interface{}) (int64, error) {
	q := new(Operation)
	q.Init(d)
	return q.insert(ptrs...)
}

func (d *Driver) Count (ptr interface{}) (int64, error) {
	q := new(Operation)
	q.Init(d)
	return q.Count(ptr)
}

func (d *Driver) Where (cmd string, args ...interface{}) *Operation {
	q := new(Operation)
	q.Init(d)
	return q.Where(cmd, args...)
}

func (d *Driver) Id (id int64) *Operation {
	q := new(Operation)
	q.Init(d)
	return q.Id(id)
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
			d.syncTable[in.Name] = in
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
	d.syncTable[in.Name] = in
	return true
}

/*
Author: ZerQAQ
 */