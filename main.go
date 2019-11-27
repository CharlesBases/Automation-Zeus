package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"CharlesBases/Automation-Zeus/template"
	"CharlesBases/Automation-Zeus/utils"
)

type Config struct {
	utils.GlobalConfig
}

var (
	db      = flag.String("d", "root:password@tcp(127.0.0.1:3306)/mysql", `database`)
	table   = flag.String("t", "", `table.multiple files are divided by ","`)
	genPath = flag.String("p", ".", `generate file path`)
	update  = flag.Bool("u", false, `update struct`)
	json    = flag.Bool("j", true, `json tag`)
	gorm    = flag.Bool("g", true, `gorm tag`)
)

func main() {
	flag.Parse()

	abspath, _ := filepath.Abs(*genPath)

	os.MkdirAll(abspath, 0755)

	config := &Config{
		GlobalConfig: utils.GlobalConfig{
			Package:     filepath.Base(abspath),
			PackagePath: abspath,
			Database: utils.Database{
				IP:     *db + "?charset=utf8mb4&sql_notes=false&sql_notes=false&timeout=90s&collation=utf8mb4_general_ci&parseTime=True&loc=Local",
				Schema: parse_schema(*db),
				Tables: make(map[string]*[]utils.TableField),
			},
			Structs: make([]*utils.Struct, 0),
			Imports: make(map[string]string, 0),
			Update:  *update,
			Json:    *json,
			Gorm:    *gorm,
		},
	}

	swg := sync.WaitGroup{}
	swg.Add(4)

	databasechannel := make(chan int64, 2)
	tableschannel := make(chan int64)

	// 连接数据库
	go func() {
		defer swg.Done()

		config.Database.InitMysql()
		databasechannel <- time.Now().UnixNano()
		databasechannel <- time.Now().UnixNano()
	}()

	// 获取数据库下所有表列表，以 order by table_name 排序
	go func() {
		defer swg.Done()

		<-databasechannel
		config.GetTable()

		tableschannel <- time.Now().UnixNano()
	}()

	// 解析已有结构体
	go func() {
		defer swg.Done()

		<-databasechannel
		config.ParseFile()
	}()

	// 是否只更新已有结构体
	go func() {
		defer swg.Done()

		<-tableschannel

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
	}()

	swg.Wait()

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

	// 生成 model
	for _, Struct := range config.Structs {
		go func(x *utils.Struct) {
			swg.Add(1)
			defer swg.Done()

			structfile := config.create(path.Join(config.PackagePath, fmt.Sprintf("%s.go", x.TableName)))
			infor := &template.Infor{Config: &config.GlobalConfig, Struct: x}
			infor.GenModel(structfile)

			structfile.Close()
		}(Struct)
	}

	swg.Wait()

	config.gofmt()

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

func (config *Config) create(filepath string) *os.File {
	os.RemoveAll(path.Join(config.PackagePath, fmt.Sprintf("%s.go", filepath)))

	fmt.Print(fmt.Sprintf(`[%s]--------`, time.Now().Format("2006-01-02 15:04:05")))
	fmt.Printf("%c[%d;%d;%dmcreate model file %s...%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 36 /*前景*/, path.Base(filepath), 0x1B)
	if file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
		fmt.Print(fmt.Sprintf(`[%s]--------`, time.Now().Format("2006-01-02 15:04:05")))
		fmt.Printf("%c[%d;%d;%dmopen file error: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B)
		os.Exit(0)
	} else {
		return file
	}
	return nil
}

func (config *Config) gofmt() {
	cmd := exec.Command("gofmt", "-l", "-w", "-s", config.PackagePath)
	err := cmd.Run()
	if err != nil {
		fmt.Print(fmt.Sprintf(`[%s]------`, time.Now().Format("2006-01-02 15:04:05")))
		fmt.Printf("%c[%d;%d;%dmgofmt error: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B)
		os.Exit(0)
	}
}
