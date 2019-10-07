package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

func (database *Database) InitMysql() {
	db, err := gorm.Open("mysql", database.IP)
	if err != nil {
		fmt.Print(fmt.Sprintf(`[%s]------`, time.Now().Format("2006-01-02 15:04:05")))
		fmt.Printf("%c[%d;%d;%dmconnection error: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B)
		os.Exit(0)
	}

	if err := db.DB().Ping(); err != nil {
		fmt.Print(fmt.Sprintf(`[%s]------`, time.Now().Format("2006-01-02 15:04:05")))
		fmt.Printf("%c[%d;%d;%dmping error: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B)
		os.Exit(0)
	}

	if DB != nil {
		DB.Close()
	}
	DB = db

	fmt.Print(fmt.Sprintf(`[%s]------`, time.Now().Format("2006-01-02 15:04:05")))
	fmt.Printf("%c[%d;%d;%dmsuccessful connection !%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 35 /*前景*/, 0x1B)
}
