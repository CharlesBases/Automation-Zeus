package utils

import "html/template"

const information_columns = "information_schema.COLUMNS"

var mysqltype = map[string]string{
	"int":        "int",
	"integer":    "int",
	"tinyint":    "int",
	"smallint":   "int",
	"mediumint":  "int",
	"bit":        "int",
	"bool":       "bool",
	"bigint":     "int64",
	"enum":       "string",
	"set":        "string",
	"varchar":    "string",
	"char":       "string",
	"tinytext":   "string",
	"mediumtext": "string",
	"text":       "string",
	"longtext":   "string",
	"blob":       "string",
	"tinyblob":   "string",
	"mediumblob": "string",
	"longblob":   "string",
	"binary":     "string",
	"varbinary":  "string",
	"json":       "string",
	"float":      "float64",
	"double":     "float64",
	"decimal":    "float64",
	"time":       "time.Time",
	"date":       "time.Time",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
}

type GlobalConfig struct {
	Package     string    // 生成结构体文件包名
	PackagePath string    // 生成结构体文件包路径
	Database    Database  //
	Structs     []*Struct // 结构体
	Update      bool
	Json        bool
	Gorm        bool
}

type Database struct {
	IP         string
	Schema     string
	Tables     map[string]*[]TableField // 实际需要输出的表
	TablesSort *[]string                // 数据库下所有表列表
}

type TableField struct {
	TableName string `gorm:"column:TABLE_NAME"`     // 表名
	Name      string `gorm:"column:COLUMN_NAME"`    // 列名
	Primary   string `gorm:"column:COLUMN_KEY"`     // 主键
	Extra     string `gorm:"column:EXTRA"`          // 自增
	IsNull    string `gorm:"column:IS_NULLABLE"`    // NOT NULL
	Type      string `gorm:"column:DATA_TYPE"`      // 类型
	Column    string `gorm:"column:COLUMN_TYPE"`    // 类型+长度
	Comment   string `gorm:"column:COLUMN_COMMENT"` // 注释
}

type Struct struct {
	StructName string
	TableName  string
	Fields     *[]StructField
}

type StructField struct {
	Name    string        // 字段
	Type    string        // 类型
	Tag     template.HTML // 标签
	Comment string        // 注释
}

// Config .
type Config struct {
	Addr    string   `toml:"addr"`
	Tables  []string `toml:"tables"`
	GenPath string   `toml:"path"`
	Update  bool     `toml:"update"`
	JsonTag bool     `toml:"json_tag"`
	GormTag bool     `toml:"gorm_tag"`
}
