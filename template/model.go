package template

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"time"
)

func (infor *Infor) GenModel(wr io.Writer) {
	temp := template.New(infor.Struct.StructName)
	temp.Funcs(template.FuncMap{
		"package": func() string {
			return infor.Config.Package
		},
		"imports": func() template.HTML {
			importsbuilder := strings.Builder{}
			for key := range infor.Config.Imports {
				importsbuilder.WriteString(fmt.Sprintf("%s\n\t", key))
			}
			return template.HTML(importsbuilder.String())
		},
		"html": func(source string) template.HTML {
			return template.HTML(source)
		},
	})
	modelTemplate, err := temp.Parse(modeltemplate)
	if err != nil {
		fmt.Print(fmt.Sprintf("[%s]----------%c[%d;%d;%dmgen model error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	}
	modelTemplate.Execute(wr, infor.Struct)
}
