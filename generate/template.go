package generate

const DefaultTemplate = `// Package {{package}} this models generated from {{.Config.Database.Schema}}. DO NOT EDIT.
package {{package}}

{{imports}}
{{range .Structs}}
// {{.StructName}} {{.Table.Comment}}
type {{.StructName}} struct {
    {{range $fieldIndex, $field := .Fields}}
	{{.Name}}   {{.Type}}   {{.Tag}}    // {{.Comment}}     {{end}}
}

func (*{{.StructName}}) TableName() string {
	return "{{.Table.Name}}"
}

{{end}}

`
