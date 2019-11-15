package utils

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"time"
)

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

func (table *{{.StructName}}) Select() error {
	return {{$orm}}.Table(table.Table()).Where("id = ? AND is_deleted = 0", table.ID).First(table).Error
}

func (table *{{.StructName}}) Update() error {
	return {{$orm}}.Table(table.Table()).Where("id = ? AND is_deleted = 0", table.ID).Updates(table).Error
}

func (table *{{.StructName}}) Delete() error {
	return {{$orm}}.Table(table.Table()).Where("id = ? AND is_deleted = 0", table.ID).Update(map[string]int{"is_deleted": 1}).Error
}

func (table *{{.StructName}}) Save(value map[string]interface{}) error {
	return {{$orm}}.Table(table.Table()).Where("id = ? AND is_deleted = 0", table.ID).Save(value).Error
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
			for orm := range config.ORM {
				importsbuilder.WriteString(fmt.Sprintf(`"%s"`, orm))
				break
			}
			return template.HTML(importsbuilder.String())
		},
		"ormcall": func() string {
			for orm := range config.ORM {
				return config.ORM[orm]
			}
			return "gorm.DB"
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
