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
	d, _ := zorm.Open("mysql", "root:123456@/test?charset=utf8")
	d.Sync(*new(user))

	t2 := user{}

	ok := d.Where("id < ?", 120).Where("name in (?)", []string{"zer", "rabbit"}).Get(&t2)
	fmt.Println(ok, t2)
}
