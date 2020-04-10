package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"orm/zorm"
)

type user struct {
	Id int64 `zorm:"pk auto_increment"`
	Name string `zorm:"not null varchar(1000)"`
}

func main()  {
	d := zorm.Open("mysql", "root:123456@/test?charset=utf8")
	d.Sync(*new(user))

	t2 := user{}
	d.Id(1).Get(&t2)
	fmt.Println(t2)

	t := make([]user, 1)
	ok := d.Where("id != ?", -1).Find(&t)
	fmt.Println(ok, t)
	fmt.Println(d.Count(new(user)))
}
