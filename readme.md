# ZORM

用Go编写的ORM框架，接口设计上大量参考了xorm。

## 如何使用？

### 连接与表同步

在引入mysqlDriver之后

```go
improt _ "github.com/go-sql-driver/mysql"
```

可以通过**Open**来打开一个数据库
```go
d, err:= zorm.Open("mysql", "userName:passWord@/dataBase?charset=utf8")
```

通过向**Sync**传入一个结构体来同步表结构，如果表不存在zorm会创建一个新表。

如果表已经存在zorm会尽量对表进行同步，若出现严重冲突（主键冲突， 列数据类型冲突等），会panic。

```go
type user struct {
	Id int64 `zorm:"pk auto_increment"`
	Name string `zorm:"not null varchar(1000)"`
	Nick string `zorm:"index"`
}

d.Sync(*new(user))
```

### 查询数据

**Get** 方法，传入一个struct指针，返回查询到的第一条数据

返回值是一个bool，表示数据是否存在
```go
ok := d.Id(20).Get(new(user))

ok := d.Where("id = ?", 10).Get(new(user))

ok := d.Where("id > ?", 120).Where("name in (?)", [2]string{"fox", "rabbit"}).Get(new(user))
//返回的是所有Where条件都被满足的数据中的第一个
```

**Find** 方法，传入一个slice指针，返回的数据条数不超过slice的长度

```go
result := make([]user, 20)

num, err := d.Where("id > ?", 20).Find(&result)
// num是找到的数据条数

num, err := d.Where("id not in (?)", []int64{10, 20}).Limit(5, 10).Find(&result)
//你也可以通过Limit来设置查找的起点和总数限制，在这个语句中，查找的约束是从第五条开始，最多不超过10条数据
```

**Count** 方法，计算数据的条数。

```go
num, err := d.Where("id < ?", 200411).Count(new(user))
```

### 插入数据

**Insert** 方法, 传入一个或者多个结构体或结构体指针。

如果传入多个结构体，任何一个结构体插入失败都会导致所有操作被回滚。

```go
num, err := d.Insert(user{Name:"owo"}, &user{Name:"QAQ"})
//num是最后一个插入的数据的Id字段（如果有的话）

num, err := d.Insert(user{Name:"rabbit"}, emotion{Content:"fox"}) 
//传入多个不同的结构体也是可以的
```

**InsertMany** 方法，传入一个slice，同时插入多条数据

```go
num, err := d.InsertMany([]user{user{Name: "=w="}, user{Name: "TvT"}})
//num是最后一个插入的数据的Id字段（如果有的话）
```

### 修改数据

**Update** 方法，传入一个结构体，根据结构体里的元素修改数据表。

```go
num, err := d.Where("id > ?", 20).Update(user{Name: ">v<"})
//num是被更新的数据条数
```

Update会自动从user结构体中提取非0和非nil得值作为需要更新的内容，因此，如果需要更新一个值为0，则此种方法将无法实现。

可以使用**Col**来指定要更新的值，这样Update除了会更新非0和非nil的值外，还会更新被指定的列。

```go
num, err := d.Where("name = ?", "0.0").Col("name").Update(user{Name:""})
//num是被更新的数据条数
```

### 删除数据

**Delete** 方法

```go
num, err := d.Where("id in (?)", []int64{1, 2, 3}).Delete(new(user))
//num是被删除的数据条数

num, err := d.Where("").Delete(new(user))
//删除所有数据
```