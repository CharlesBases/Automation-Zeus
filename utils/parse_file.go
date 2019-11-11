package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func (config *GlobalConfig) ParseFile() {
	fmt.Print(fmt.Sprintf(`[%s]--------`, time.Now().Format("2006-01-02 15:04:05")))
	fmt.Printf("%c[%d;%d;%dmparse package: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 36 /*前景*/, config.PackagePath, 0x1B)
	filename, err := ioutil.ReadDir(config.PackagePath)
	if err != nil {
		fmt.Print(fmt.Sprintf(`[%s]--------`, time.Now().Format("2006-01-02 15:04:05")))
		fmt.Printf("%c[%d;%d;%dmparse package error: %s%c[0m\n", 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B)
	}

	for _, val := range filename {
		if strings.HasSuffix(val.Name(), ".go") {
			config.Database.Tables[ensnake(strings.TrimRight(val.Name(), ".go"))] = &[]TableField{}
		}
	}
}
