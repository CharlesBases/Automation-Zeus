package utils

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"time"
)

const modeltemplate = `{{package}}

type {{.StructName}} struct {                               {{range $fieldIndex, $field := .Fields}}
	{{.Name}}   {{.Type}}   {{.Tag}}    // {{.Comment}}     {{end}}
}

func (*{{.StructName}}) Table() string {
	return "{{tablename .StructName}}"
}

func(table *{{.StructName}}) Insert() error {
	return orm.Gorm().Table(table.Table()).Create(table).Error
}

func(table *{{.StructName}}) Select(ID int64) error {
	return orm.Gorm().Table(table.Table()).Where("id = ? AND is_deleted = 0", ID).First(table).Error
}

func(table *{{.StructName}}) Update(ID int64) error {
	return orm.Gorm().Table(table.Table()).Where("id = ? AND is_deleted = 0", ID).Update(table).Limit(1).Error
}

func(table *{{.StructName}}) Delete(ID int64) error {
	return orm.Gorm().Table(table.Table()).Where("id = ? AND is_deleted = 0", ID).Updates(table).Limit(1).Error
}

`

func (config *GlobalConfig) GenModel(Struct *Struct, wr io.Writer) {
	temp := template.New(Struct.StructName)
	temp.Funcs(template.FuncMap{
		"package": func() template.HTML {
			builder := strings.Builder{}
			// package
			builder.WriteString(fmt.Sprintf("package %s\n\n", config.Package))
			// import
			if config.Import != nil {
				builder.WriteString(fmt.Sprintf("import (\n%s)", func() string {
					importBuilder := strings.Builder{}
					for key, val := range config.Import {
						importBuilder.WriteString(fmt.Sprintf("\t%s %s\n", val, key))
					}

					importBuilder.WriteString("\t")
					importBuilder.WriteString(fmt.Sprintf(`"gitlab.ifchange.com/bot/gokitcommon/orm"`))
					importBuilder.WriteString("\n")
					return importBuilder.String()
				}()))
			}
			return template.HTML(builder.String())
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
