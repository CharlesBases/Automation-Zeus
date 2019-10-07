package Automation_zeus

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/CharlesBases/Automation-zeus/utils"
)

type Config struct {
	utils.GlobalConfig
}

var (
	db     = flag.String("d", "root:password@tcp(127.0.0.1:3306)/mysql", `database`)
	table  = flag.String("t", "", `table.multiple files are divided by ","`)
	file   = flag.String("f", "./model.go", `generate model file name`)
	update = flag.Bool("u", false, `update struct`)
	json   = flag.Bool("j", true, `json tag`)
	gorm   = flag.Bool("g", true, `gorm tag`)
)

func main() {
	flag.Parse()

	abspath, _ := filepath.Abs(*file)

	config := &Config{
		GlobalConfig: utils.GlobalConfig{
			Package:  filepath.Base(filepath.Dir(abspath)),
			Filepath: abspath,
			Variable: map[string]string{"schema": fmt.Sprintf(`"%s"`, parse_schema(*db))},
			Import:   make(map[string]string, 0),
			Database: &utils.Database{
				IP:     *db + "?charset=utf8mb4&parseTime=True&loc=Local",
				Schema: parse_schema(*db),
				Tables: make(map[string]*[]utils.TableField),
			},
			Structs: &[]utils.Struct{},
			Update:  *update,
			Json:    *json,
			Gorm:    *gorm,
		},
	}

	// 连接数据库
	config.Database.InitMysql()

	// 获取数据库下所有表列表，以 order by table_name 排序
	config.GetTable()

	// 解析已有结构体
	if astFile, err := parser.ParseFile(token.NewFileSet(), config.Filepath, nil, parser.ParseComments); err == nil {
		config.ParseFile(astFile)
	}

	// 是否只更新已有结构体
	if !config.Update {
		if len(*table) != 0 {
			for _, val := range strings.Split(*table, ",") {
				config.Database.Tables[val] = &[]utils.TableField{}
			}
		} else {
			for _, val := range *(config.Database.TablesSort) {
				config.Database.Tables[val] = &[]utils.TableField{}
			}
		}
	}

	fmt.Print(fmt.Sprintf(`[%s]--------`, time.Now().Format("2006-01-02 15:04:05")))
	fmt.Printf("%c[%d;%d;%dmparse database: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 36 /*前景*/, config.Database.Schema, 0x1B)

	// 获取表结构
	for _, val := range *config.Database.TablesSort {
		if _, ok := config.Database.Tables[val]; ok {
			fmt.Println(fmt.Sprintf(`[%s]----------find table: %s`, time.Now().Format("2006-01-02 15:04:05"), val))
			config.Database.GetTable(val)
			config.ParseTable(config.Database.Tables[val])
		}
	}

	// 生成文件
	modelfile := config.createModel()
	config.GenModel(modelfile)
	modelfile.Close()

	config.gofmt(config.Filepath)

	fmt.Print(fmt.Sprintf(`[%s]------`, time.Now().Format("2006-01-02 15:04:05")))
	fmt.Printf("%c[%d;%d;%dmcomplete !%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 35 /*前景*/, 0x1B)
}

// 获取数据库名
func parse_schema(ip string) (schema string) {
	start := strings.LastIndex(ip, "/")
	if start != -1 {
		schema = ip[start+1:]
	}
	return
}

func (config *Config) createModel() *os.File {
	os.RemoveAll(config.Filepath)

	fmt.Print(fmt.Sprintf(`[%s]--------`, time.Now().Format("2006-01-02 15:04:05")))
	fmt.Printf("%c[%d;%d;%dmcreate model file...%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 36 /*前景*/, 0x1B)
	if file, err := os.OpenFile(config.Filepath, os.O_RDWR|os.O_CREATE, 0755); err != nil {
		fmt.Print(fmt.Sprintf(`[%s]--------`, time.Now().Format("2006-01-02 15:04:05")))
		fmt.Printf("%c[%d;%d;%dmopen file error: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B)
		os.Exit(0)
	} else {
		return file
	}
	return nil
}

func (config *Config) gofmt(filepath string) {
	cmd := exec.Command("gofmt", "-l", "-w", "-s", filepath)
	err := cmd.Run()
	if err != nil {
		fmt.Print(fmt.Sprintf(`[%s]------`, time.Now().Format("2006-01-02 15:04:05")))
		fmt.Printf("%c[%d;%d;%dmgofmt error: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B)
		os.Exit(0)
	}
}
