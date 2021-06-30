package generate

const modeltemplate = `// Package {{package}} this model is generate from schema {{.Config.Database.Schema}}
package {{package}}

import (
	{{imports}}
	"github.com/jinzhu/gorm"
)

{{range .Structs}}
// {{.StructName}} {{.Table.Comment}}
type {{.StructName}} struct {
	{{gormDB}}														
															{{range $fieldIndex, $field := .Fields}}
	{{.Name}}   {{.Type}}   {{.Tag}}    // {{.Comment}}     {{end}}
}

func (*{{.StructName}}) TableName() string {
	return "{{.Table.Name}}"
}

{{end}}

`
