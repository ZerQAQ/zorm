package zorm

import (
	"database/sql"
	"github.com/ZerQAQ/zorm/set"
	"github.com/ZerQAQ/zorm/table"
)

func (d *Driver) Connect (name string, sour string) error {
	var err error
	d.Database, err = sql.Open(name, sour)
	return err
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

		var stringVal = string(val.([]byte))
		d.tableSet.Insert(stringVal)
	}
	d.syncTable = make(map[string]table.Table)
}