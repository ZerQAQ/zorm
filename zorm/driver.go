package zorm

import (
	"database/sql"
	"errors"
	"github.com/ZerQAQ/zorm/global"
	"github.com/ZerQAQ/zorm/set"
	"github.com/ZerQAQ/zorm/table"
	"reflect"
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
	t := global.UnpackPtr(reflect.ValueOf(s)).Type()

	if t.Kind() != reflect.Struct {panic(errors.New("zorm: the parameter send to Sync must be a struct or struct pointer"))}

	tableStruct := table.ParseStruct(t)
	if d.tableSet.Find(tableStruct.Name) {
		databaseTable := d.makeTable(tableStruct.Name)

		if !table.IsContain(&databaseTable, &tableStruct){
			panic(errors.New("zorm: conflict happen cant sync"))
		}

		if table.IsContain(&tableStruct, &databaseTable){
			d.syncTable[tableStruct.Name] = tableStruct
			return true
		}

		addition := table.Sub(&tableStruct, &databaseTable)
		tableStruct.RowsName = databaseTable.RowsName

		tableStruct.UpdateRowList(&addition)
		d.alterTable(&addition)
		//debug.PrtTable(&ad)
	} else {
		//fmt.Println("Creating")
		d.createTable(&tableStruct)
		//debug.Prt(&in)
	}
	d.syncTable[tableStruct.Name] = tableStruct
	return true
}

/*
Author: ZerQAQ
 */