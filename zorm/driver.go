package zorm

import (
	"database/sql"
	"orm/debug"
	"orm/set"
	"orm/table"
)

type Driver struct {
	db *sql.DB
	tableSet *set.Set
}

func Open(name string, sour string) *Driver {
	ret := new(Driver)
	ret.Connect(name, sour)
	ret.initTables()
	return ret
}

func (d *Driver) Sync (s interface{}) bool {
	in := table.ParseStruct(s)
	if d.tableSet.Find(in.Name) {
		tb := d.makeTable(in.Name)

		debug.PrtTable(&tb)
		debug.PrtTable(&in)

		if !table.IsContain(&tb, &in){
			panic("Conflict, cant sync.")
		}

		if table.IsContain(&in, &tb){
			return true
		}

		ad := table.Sub(&in, &tb)
		debug.PrtTable(&ad)
		return true
	}
	return true
}