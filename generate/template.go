package generate

const DefaultTemplate = `// Package {{package}} this models is generated from {{.Config.Database.Schema}}
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
