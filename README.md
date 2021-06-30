# Automation

代码自动化

# mysql-gen-go
根据MySql表结构生成Golang结构体，支持json标签，gorm标签。需要使用gofmt代码格式化工具。

## 参数说明
```sh
	-f string
	    config file (default "./reverse.toml")
```

## 配置文件说明
```sh
	addr: 数据库地址。(例: root:password@tcp(127.0.0.1:3306)/mysql)
	path: models 生成路径。(例: ./models)
```

## 生成示例
```go
package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

var UsersModel = new(Users)

func init() {
	UsersModel.db = DatabaseModel.Table(UsersModel.Table())
}

// Users 用户
type Users struct {
	db *gorm.DB `json:"-" gorm:"-"`

	ID              uint      `json:"id" gorm:"column:id;type:int(10) unsigned;not null;primary_key;auto_increment"` // 用户编号
	Name            string    `json:"name" gorm:"column:name;type:varchar(40);not null"`                             // 用户名
	Pwd             string    `json:"pwd" gorm:"column:pwd;type:varchar(200);not null"`                              // 密码
	Birthday        time.Time `json:"birthday" gorm:"column:birthday;type:date;not null"`                            // 生日
}

func (*Users) TableName() string {
	return "users"
}

```
