package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"CharlesBases/mysql-gen-go/generate"
	"CharlesBases/mysql-gen-go/utils"
)

type Config struct {
	utils.GlobalConfig
}

var (
	file = flag.String("f", "./reverse.toml", "config file")
)

var toml *utils.Config

func main() {
	flag.Parse()

	toml = utils.ParseConfig(*file)

	abspath, _ := filepath.Abs(toml.GenPath)

	os.MkdirAll(abspath, 0755)

	config := &Config{
		GlobalConfig: utils.GlobalConfig{
			Package:     filepath.Base(abspath),
			PackagePath: abspath,
			Database: utils.Database{
				IP:     toml.Addr + "?charset=utf8mb4&sql_notes=false&sql_notes=false&timeout=90s&collation=utf8mb4_general_ci&parseTime=True&loc=Local",
				Schema: parse_schema(toml.Addr),
				Tables: make(map[string]*[]utils.TableField),
			},
			Structs: make([]*utils.Struct, 0),
			Imports: make(map[string]string, 0),
		},
	}

	// 连接数据库
	config.Database.InitMysql()

	// 获取数据库下所有表列表，以 order by table_name 排序
	config.GetTables()

	fmt.Print(fmt.Sprintf("[%s]--------%c[%d;%d;%dmparse database: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 36 /*前景*/, config.Database.Schema, 0x1B))

	// 获取表结构
	for _, val := range *config.Database.TablesSort {
		if _, ok := config.Database.Tables[val]; ok {
			fmt.Println(fmt.Sprintf(`[%s]----------find table: %s`, time.Now().Format("2006-01-02 15:04:05"), val))
			config.Database.GetTableFields(val)
			config.ParseTable(config.Database.Tables[val])
		}
	}

	// 生成 model
	structfile := config.create(path.Join(config.PackagePath, "models.go"))
	infor := &generate.Infor{Config: &config.GlobalConfig, Structs: config.Structs}
	infor.GenModel(structfile)
	structfile.Close()

	config.gofmt()

	fmt.Print(fmt.Sprintf("[%s]------%c[%d;%d;%dmcomplete !%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 35 /*前景*/, 0x1B))
}

// 获取数据库名
func parse_schema(ip string) (schema string) {
	start := strings.LastIndex(ip, "/")
	if start != -1 {
		schema = ip[start+1:]
	}
	return
}

func (config *Config) create(filepath string) *os.File {
	os.RemoveAll(path.Join(config.PackagePath, fmt.Sprintf("%s.go", filepath)))

	fmt.Print(fmt.Sprintf("[%s]--------%c[%d;%d;%dmcreate model file %s...%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 36 /*前景*/, path.Base(filepath), 0x1B))
	if file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
		fmt.Print(fmt.Sprintf("[%s]--------%c[%d;%d;%dmopen file error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	} else {
		return file
	}
	return nil
}

func (config *Config) gofmt() {
	cmd := exec.Command("gofmt", "-l", "-w", "-s", config.PackagePath)
	err := cmd.Run()
	if err != nil {
		fmt.Print(fmt.Sprintf("[%s]------%c[%d;%d;%dmgofmt error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	}
}
