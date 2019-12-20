package template

const modeltemplate = `// this model is generate for table {{.TableName}}
package {{package}}

import (
	"fmt"
	"sync"
	{{imports}}
	"github.com/jinzhu/gorm"
	"gitlab.ifchange.com/bot/gokitcommon/log"
	"gitlab.ifchange.com/bot/gokitcommon/orm"
)

type {{.StructName}} struct {                               {{range $fieldIndex, $field := .Fields}}
	{{.Name}}   {{.Type}}   {{.Tag}}    // {{.Comment}}     {{end}}
}

func New{{.StructName}}() *{{.StructName}} {
    return new({{.StructName}})
}

func (*{{.StructName}}) Table() string {
	return "{{.TableName}}"
}

func (table *{{.StructName}}) DB() *gorm.DB {
	return orm.Gorm().Table(table.Table())
}

func (table *{{.StructName}}) Insert() error {
	return table.DB().Create(table).Error
}

func (table *{{.StructName}}) Select() error {
	return table.DB().Where("id = ? AND is_deleted = 0", table.ID).First(table).Error
}

func (table *{{.StructName}}) Update() error {
	return table.DB().Where("id = ? AND is_deleted = 0", table.ID).Updates(table).Error
}

func (table *{{.StructName}}) Delete() error {
	return table.DB().Where("id = ? AND is_deleted = 0", table.ID).Update(map[string]int{"is_deleted": 1}).Error
}

func (table *{{.StructName}}) First(query interface{}, args ...interface{}) error {
	return table.DB().Where(query, args...).First(table).Error
}

func (table *{{.StructName}}) Selects(query interface{}, args ...interface{}) (*[]{{.StructName}}, error) {
	list := make([]{{.StructName}}, 0)
	err := table.DB().Where(query, args...).Find(list).Error
	return &list, err
}

func (table *{{.StructName}}) Updates(datas map[string]interface{}, query interface{}, args ...interface{}) error {
	return table.DB().Where(query, args...).Updates(datas).Error
}

func (table *{{.StructName}}) Inserts(tables *[]{{.StructName}}) error {
	swg := sync.WaitGroup{}
	swg.Add(len(*tables))
	errorchannel := make(chan error, len(*tables))
	tx := table.DB().Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	for _, table := range *tables {
		go func(x *{{.StructName}}) {
			defer swg.Done()
			if err := tx.Create(x).Error; err != nil {
				errorchannel {{html "<-"}} err
				log.Errorf(fmt.Sprintf("Inserts %s error: %v", x.Table(), err))
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

`
