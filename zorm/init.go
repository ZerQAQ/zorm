package zorm

import (
	"database/sql"
	"orm/set"
	"orm/table"
)

func (d *Driver) Connect (name string, sour string)  {
	var err error
	d.Database, err = sql.Open(name, sour)
	if err != nil {panic(err)}
}

func (d *Driver) init (){
	d.tableSet = set.MakeSet()
	rows, err := d.Database.Query("show tables")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var val interface{}

		err := rows.Scan(&val)
		if err != nil {panic(err)}

		var strval = string(val.([]byte))
		d.tableSet.Insert(strval)
	}
	d.SyncTable = make(map[string]table.Table)
}