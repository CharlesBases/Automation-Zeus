package template

const modeltemplate = `// this model is generate for table {{.TableName}}
package {{package}}

import (
	"fmt"
	"sync"
	{{imports}}
	"github.com/jinzhu/gorm"
	"github.com/CharlesBases/common/log"
)

var {{.StructName}}Model = new({{.StructName}})

func init() {
	{{.StructName}}Model.db = DatabaseModel.Table({{.StructName}}Model.Table())
}

type {{.StructName}} struct {
	{{gormDB}}														
															{{range $fieldIndex, $field := .Fields}}
	{{.Name}}   {{.Type}}   {{.Tag}}    // {{.Comment}}     {{end}}
}

func (*{{.StructName}}) Table() string {
	return "{{.TableName}}"
}

func (*{{.StructName}}) Insert(table *Interviews) error {
	return {{.StructName}}Model.db.Create(table).Error
}

func (*{{.StructName}}) Deletes(query interface{}, args ...interface{}) error {
	return {{.StructName}}Model.db.Where(query, args...).Update(map[string]int{"is_deleted": 1}).Error
}

func (*{{.StructName}}) Updates(params map[string]interface{}, query interface{}, args ...interface{}) error {
	return {{.StructName}}Model.db.Where(query, args...).Updates(params).Error
}

func (*{{.StructName}}) First(query interface{}, args ...interface{}) (*{{.StructName}}, error) {
	result := new({{.StructName}})
	err := {{.StructName}}Model.db.Where(query, args...).First(result).Error
	return result, err
}

func (*{{.StructName}}) Selects(query interface{}, args ...interface{}) (*[]{{.StructName}}, error) {
	list := make([]{{.StructName}}, 0)
	err := {{.StructName}}Model.db.Where(query, args...).Find(list).Error
	return &list, err
}

func (*{{.StructName}}) Inserts(tables *[]{{.StructName}}) error {
	swg := sync.WaitGroup{}
	swg.Add(len(*tables))
	errorchannel := make(chan error, len(*tables))
	tx := {{.StructName}}Model.db.Begin()
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
