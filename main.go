package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"mysql-gen-go/generate"
	"mysql-gen-go/logger"
	"mysql-gen-go/utils"
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
			Database: &utils.Database{
				IP:     toml.Addr + "?charset=utf8mb4&sql_notes=false&sql_notes=false&timeout=60s&collation=utf8mb4_general_ci&parseTime=True&loc=Local",
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

	logger.Infor("parse database: %s", config.Database.Schema)

	// 获取表结构
	for _, table := range *config.Database.TablesSort {
		if _, ok := config.Database.Tables[table]; ok {
			logger.Debug("find table: %s", table)
			config.Database.GetTableFields(table)
			config.ParseTable(config.Database.Tables[table])
		}
	}

	// 生成 model
	structfile := config.create(path.Join(config.PackagePath, "models.go"))
	infor := &generate.Infor{
		Config:  &config.GlobalConfig,
		Structs: config.Structs,
		Template: func() string {
			if toml.Template != "" {
				return toml.Template
			}
			return generate.DefaultTemplate
		}(),
	}
	infor.GenModel(structfile)
	structfile.Close()

	config.gofmt()

	logger.Infor("complete !")
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

	logger.Infor("creating file[%s]...", path.Base(filepath))
	if file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
		logger.Error("open file error: %v", err)
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
		logger.Error("gofmt error: %v", err)
		os.Exit(1)
	}
}
