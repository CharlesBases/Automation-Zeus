package utils

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const modeltemplate = `{{package}}
{{range $structIndex, $struct := .Structs}}
type {{.StructName}} struct {                               {{range $fieldIndex, $field := .Fields}}
	{{.Name}}   {{.Type}}   {{.Tag}}    // {{.Comment}}     {{end}}
}

func (*{{.StructName}}) Table() string {
	return schema + ".{{tablename .StructName}}"
}
{{end}}
`

func (config *GlobalConfig) GenModel(wr io.Writer) {
	temp := template.New(filepath.Base(config.Filepath))
	temp.Funcs(template.FuncMap{
		"package": func() template.HTML {
			builder := strings.Builder{}
			// package
			builder.WriteString(fmt.Sprintf("package %s\n\n", config.Package))
			// import
			if config.Import != nil {
				builder.WriteString(fmt.Sprintf("import (\n%s)\n\n", func() string {
					importBuilder := strings.Builder{}
					for key, val := range config.Import {
						importBuilder.WriteString(fmt.Sprintf("\t%s %s\n", val, key))
					}
					return importBuilder.String()
				}()))
			}
			// var
			builder.WriteString(fmt.Sprintf("var (\n%s)", func() string {
				varBuilder := strings.Builder{}
				for key, val := range config.Variable {
					varBuilder.WriteString(fmt.Sprintf("\t%s = %s\n", key, val))
				}
				return varBuilder.String()
			}()))
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
	modelTemplate.Execute(wr, config)
}
