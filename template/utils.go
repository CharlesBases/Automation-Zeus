package template

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"time"

	"CharlesBases/Automation-Zeus/utils"
)

type Infor struct {
	Config *utils.GlobalConfig
	Struct *utils.Struct
}

func (infor *Infor) GenUtils(wr io.Writer) {
	temp := template.New("utils")
	temp.Funcs(template.FuncMap{
		"package": func() string {
			return infor.Config.Package
		},
	})
	utilsTemplate, err := temp.Parse(utilsTemplate)
	if err != nil {
		fmt.Print(fmt.Sprintf("[%s]----------%c[%d;%d;%dmgen utils error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	}
	utilsTemplate.Execute(wr, infor.Config)
}
