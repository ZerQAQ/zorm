package zorm

import (
	"database/sql"
	"orm/set"
)

func (d *Driver) Connect (name string, sour string)  {
	var err error
	d.db, err = sql.Open(name, sour)
	if err != nil {panic(err)}
	d.initTables()
}

func (d *Driver) initTables (){
	d.tableSet = set.MakeSet()
	rows, err := d.db.Query("show tables")
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
}