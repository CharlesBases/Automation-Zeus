package utils

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"time"
)

// 自封装的 gorm 路径
const ormimports = "github.com/CharlesBases/common/orm/gorm"

// gorm 调用
var romcall = map[string]string{
	ormimports: "gorm.DB",
}

const modeltemplate = `// this model is generate for {{.StructName}} {{$orm:=ormcall}}
package {{package}}

import (
{{imports}}
)

type {{.StructName}} struct {                               {{range $fieldIndex, $field := .Fields}}
	{{.Name}}   {{.Type}}   {{.Tag}}    // {{.Comment}}     {{end}}
}

func (*{{.StructName}}) Table() string {
	return "{{tablename .StructName}}"
}

func (table *{{.StructName}}) Insert() error {
	return {{$orm}}.Table(table.Table()).Create(table).Error
}

func (table *{{.StructName}}) Select(ID int64) error {
	return {{$orm}}.Table(table.Table()).Where("id = ? AND is_deleted = 0", ID).First(table).Error
}

func (table *{{.StructName}}) Update(ID int64) error {
	return {{$orm}}.Table(table.Table()).Where("id = ? AND is_deleted = 0", ID).Updates(table).Error
}

func (table *{{.StructName}}) Delete(ID ...int64) error {
	return {{$orm}}.Table(table.Table()).Where("id IN (?) AND is_deleted = 0", ID).Update(map[string]int{"is_deleted": 1}).Error
}

func (table *{{.StructName}}) Save(ID int64, value map[string]interface{}) error {
	return {{$orm}}.Table(table.Table()).Where("id = ? AND is_deleted = 0", ID).Save(value).Error
}

`

func (config *GlobalConfig) GenModel(Struct *Struct, wr io.Writer) {
	temp := template.New(Struct.StructName)
	temp.Funcs(template.FuncMap{
		"package": func() string {
			return config.Package
		},
		"imports": func() template.HTML {
			importsbuilder := strings.Builder{}
			// others import
			for key, val := range config.Import {
				importsbuilder.WriteString(fmt.Sprintf("\t%s %s\n\t", val, key))
			}
			// gorm import
			importsbuilder.WriteString(fmt.Sprintf(`"%s"`, ormimports))
			return template.HTML(importsbuilder.String())
		},
		"ormcall": func() string {
			return romcall[ormimports]
		},
		"tablename": ensnake,
	})
	modelTemplate, err := temp.Parse(modeltemplate)
	if err != nil {
		fmt.Print(fmt.Sprintf(`[%s]----------`, time.Now().Format("2006-01-02 15:04:05")))
		fmt.Printf("%c[%d;%d;%dmgen model error: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B)
		os.Exit(1)
	}
	modelTemplate.Execute(wr, Struct)
}
