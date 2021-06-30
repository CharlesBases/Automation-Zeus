package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

// InitMysql 连接数据库
func (database *Database) InitMysql() {
	db, err := gorm.Open("mysql", database.IP)
	if err != nil {
		fmt.Print(fmt.Sprintf("[%s]------%c[%d;%d;%dmconnection error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	}

	if err := db.DB().Ping(); err != nil {
		fmt.Print(fmt.Sprintf("[%s]------%c[%d;%d;%dmping error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	}

	if DB != nil {
		DB.Close()
	}
	DB = db

	fmt.Print(fmt.Sprintf("[%s]------%c[%d;%d;%dmsuccessful connection !%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 35 /*前景*/, 0x1B))
}

// GetTables 获取当前库下所有表名
func (config *GlobalConfig) GetTables() {
	tables := make([]string, 0)
	err := DB.Table(information_columns).
		Where("TABLE_SCHEMA = ?", config.Database.Schema).
		Group("TABLE_NAME").
		Order("TABLE_NAME").
		Pluck("TABLE_NAME", &tables).
		Error
	if err != nil {
		fmt.Print(fmt.Sprintf("[%s]--------%c[%d;%d;%dminformation_columns error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	}
	config.Database.TablesSort = &tables

	for _, table := range tables {
		config.Database.Tables[table] = &[]TableField{}
	}
}

// GetTable 获取表信息
func (database *Database) GetTable(tableName string) *Table {
	table := new(Table)

	err := DB.Table(information_tables).
		Select([]string{
			"TABLE_NAME",
			"TABLE_COMMENT",
		}).
		Where("TABLE_SCHEMA = ? AND TABLE_NAME = ?", database.Schema, tableName).
		First(table).
		Error
	if err != nil {
		fmt.Print(fmt.Sprintf("[%s]--------%c[%d;%d;%dminformation_tables error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	}

	table.Comment = strings.TrimSpace(table.Comment)
	return table
}

// GetTableFields 获取表字段
func (database *Database) GetTableFields(tableName string) {
	err := DB.Table(information_columns).
		Select([]string{
			"TABLE_NAME",
			"COLUMN_NAME",
			"COLUMN_KEY",
			"EXTRA",
			"IS_NULLABLE",
			"DATA_TYPE",
			"COLUMN_TYPE",
			"COLUMN_COMMENT",
		}).
		Where("TABLE_SCHEMA = ? AND TABLE_NAME = ?", database.Schema, tableName).
		Order("ORDINAL_POSITION").
		Find(database.Tables[tableName]).
		Error
	if err != nil {
		fmt.Print(fmt.Sprintf("[%s]--------%c[%d;%d;%dminformation_columns error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	}
}
