package generate

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"mysql-gen-go/logger"
	"mysql-gen-go/utils"
)

type Infor struct {
	Config   *utils.GlobalConfig
	Structs  []*utils.Struct
	Template string
}

func (infor *Infor) GenModel(wr io.Writer) {
	temp := template.New(infor.Config.Database.Schema)
	temp.Funcs(template.FuncMap{
		"package": func() string {
			return infor.Config.Package
		},
		"imports": func() template.HTML {
			if len(infor.Config.Imports) != 0 {
				importsbuilder := strings.Builder{}
				importsbuilder.WriteString("import (\n")
				for key := range infor.Config.Imports {
					importsbuilder.WriteString(fmt.Sprintf("\t%s\n", key))
				}
				importsbuilder.WriteString(")")
				return template.HTML(importsbuilder.String())
			}
			return ""
		},
		"html": func(source string) template.HTML {
			return template.HTML(source)
		},
		"gormDB": func() template.HTML {
			strBuilder := strings.Builder{}
			strBuilder.WriteString("db *gorm.DB " + "`" + `json:"-" gorm:"-"` + "`")
			return template.HTML(strBuilder.String())
		},
	})
	modelTemplate, err := temp.Parse(infor.Template)
	if err != nil {
		logger.Error("generate model error: %v", err)
		os.Exit(1)
	}
	modelTemplate.Execute(wr, infor)
}
