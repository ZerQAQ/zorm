package main

import (
	_ "github.com/go-sql-driver/mysql"
	"orm/zorm"
)

type user struct {
	Id int64 `zorm:"pk unique auto_increment"`
	Name string `zorm:"unique"`
	Score int64
}

func main()  {
	d := zorm.Open("mysql", "root:123456@/test?charset=utf8")
	d.Sync(*new(user))
}
