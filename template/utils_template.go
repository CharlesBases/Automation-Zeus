package template

const utilsTemplate = `
package {{package}}

import (
	orm "github.com/CharlesBases/common/orm/gorm"
)

var (
	DatabaseModel = orm.Gorm()
)

`
