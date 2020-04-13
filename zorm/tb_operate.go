package zorm

import (
	"github.com/ZerQAQ/zorm/table"
	"strings"
)

func (d *Driver) alterTable (t *table.Table)  {
	for _, elm := range t.Rows{
		sql := "alter table " + t.Name +
			" add column " + elm.Name + " " + elm.Type
		if elm.Pk {sql += " primary key "}
		if !elm.Null {sql += " not null "}
		if elm.AutoIncrement {sql += " auto_increment "}
		//fmt.Println(sql)
		d.Database.Exec(sql)
	}
	for _, elm := range t.Indexs{
		sql := "create"
		if elm.Unique {sql += " unique index "} else{
			sql += " index "
		}
		sql += elm.Name + " on " + t.Name + "(" + elm.ColName + ")"
		//fmt.Println(sql)
		d.Database.Exec(sql)
	}
}

func (d *Driver) createTable (t *table.Table)  {
	sql := "create table " + t.Name + "("
	notFir := false
	for _, elm := range t.Rows{
		if notFir {sql += ",\n"}
		notFir = true
		sql += elm.Name + " " + elm.Type + " "
		if elm.Pk {sql += " primary key "}
		if !elm.Null {sql += " not null "}
		if elm.AutoIncrement {sql += " auto_increment "}
	}
	for _, elm := range t.Indexs{
		if strings.Contains(elm.Name, "pk_") {continue}
		sql += ",\n"
		if elm.Unique {sql += "unique "}
		sql += "index(" + elm.ColName + ")"
	}
	sql += ")"

	d.Database.Exec(sql)
}
