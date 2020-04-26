# Automation

代码自动化

## 导航
* [Zeus](https://github.com/CharlesBases/Automation-Zeus) - - - - - - - - - - - - - - - - - - - - - - MySQL 表结构自动生成 Golang 结构体
* [Hera](https://github.com/CharlesBases/Automation-Hera) - - - - - - - - - - - - - - - - - - - - - - Golang 结构体生成 Json 文档
* [Hephaestus](https://github.com/CharlesBases/Automation-Hephaestus)
* [Poseidon](https://github.com/CharlesBases/Automation-Poseidon) - - - - - - - - - - - - - - - - - - - 根据接口生成 RPC 文件
* [Hestia](https://github.com/CharlesBases/Automation-Hestia)
* [Hades](https://github.com/CharlesBases/Automation-Hades)
* [Hermes](https://github.com/CharlesBases/Automation-Hermes)
* [Ares](https://github.com/CharlesBases/Automation-Ares)
* [Artemis](https://github.com/CharlesBases/Automation-Artemis)
* [Apollo](https://github.com/CharlesBases/Automation-Apollo)
* [Aphrodite](https://github.com/CharlesBases/Automation-Aphrodite)
* [Athena](https://github.com/CharlesBases/Automation-Athena)

# Automation-Zeus
根据MySql表结构生成Golang结构体，支持同一个库下多表创建，支持json标签，gorm标签，支持文件追加，支持结构体更新。需要使用gofmt代码格式化工具。

## 参数说明
```sh
	-d string
	    database (default "root:password@tcp(127.0.0.1:3306)/mysql")
	-f string
	    generate model file name (default "./model.go")
	-t string
	    table.multiple files are divided by ","
	-g    gorm tag (default true)
	-j    json tag (default true)
	-u    update struct (default false)
```

## 使用说明
```sh
	-d 使用时需要加引号
	-u 使用 -u 参数时，只更新已有的结构体
	-t 指定表名时，生成指定表，并更新已有结构体；未指定表名时，生成数据库下所以表，并更新已有结构体
	-结构体以 order by table_name 排序
```

## 生成示例
```go
package model

import (
	"sync"
	"time"

	"github.com/CharlesBases/common/log"
	"github.com/jinzhu/gorm"
)

var UsersModel = new(Users)

func init() {
	UsersModel.db = DatabaseModel.Table(UsersModel.Table())
}

type Users struct {
	db *gorm.DB `json:"-" gorm:"-"`

	ID              uint      `json:"id" gorm:"column:id;type:int(10) unsigned;not null;primary_key;auto_increment"` // 用户编号
	Name            string    `json:"name" gorm:"column:name;type:varchar(40);not null"`                             // 用户名
	Pwd             string    `json:"pwd" gorm:"column:pwd;type:varchar(200);not null"`                              // 密码
	Birthday        time.Time `json:"birthday" gorm:"column:birthday;type:date;not null"`                            // 生日
}

func (*Users) Table() string {
	return "users"
}

func (*Users) Insert(table *Users) error {
	return UsersModel.db.Create(table).Error
}

func (*Users) Deletes(query interface{}, args ...interface{}) error {
	return UsersModel.db.Where(query, args...).Update(map[string]int{"is_deleted": 1}).Error
}

func (*Users) Updates(params map[string]interface{}, query interface{}, args ...interface{}) error {
	return UsersModel.db.Where(query, args...).Updates(params).Error
}

func (*Users) Pluck(resultPoint interface{}, column string, query interface{}, args ...interface{}) error {
	return UsersModel.db.Where(query, args...).Pluck(column, resultPoint).Error
}

func (*Users) First(query interface{}, args ...interface{}) (error, *Users) {
	result := new(Users)
	err := UsersModel.db.Where(query, args...).First(result).Error
	return err, result
}

func (*Users) Selects(query interface{}, args ...interface{}) (error, *[]Users) {
	list := make([]Users, 0)
	err := UsersModel.db.Where(query, args...).Find(&list).Error
	return err, &list
}

func (*Users) IsExist(query interface{}, args ...interface{}) (error, bool) {
	var (
		result  = new(Users)
		isExist bool
	)
	err := UsersModel.db.Where(query, args...).First(result).Error
	if result != nil {
		isExist = true
	}
	return err, isExist
}

func (*Users) Inserts(tables *[]Users) error {
	swg := sync.WaitGroup{}
	swg.Add(len(*tables))
	errorchannel := make(chan error, len(*tables))
	tx := UsersModel.db.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	for _, table := range *tables {
		go func(x *Users) {
			defer swg.Done()
			if err := tx.Create(x).Error; err != nil {
				errorchannel <- err
				log.Errorf("Inserts %s error: %v", x.Table(), err)
			}
		}(&table)
	}
	swg.Wait()
	close(errorchannel)
	for err := range errorchannel {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

```
