package generate

const modeltemplate = `// Package {{package}} this model is generate from table {{.TableName}}
package {{package}}

import (
	{{imports}}
	"github.com/jinzhu/gorm"
)

type {{.StructName}} struct {
	{{gormDB}}														
															{{range $fieldIndex, $field := .Fields}}
	{{.Name}}   {{.Type}}   {{.Tag}}    // {{.Comment}}     {{end}}
}

func (*{{.StructName}}) Table() string {
	return "{{.TableName}}"
}

`
