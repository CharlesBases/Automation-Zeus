package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

var once sync.Once

// ParseConfig 解析配置文件
func ParseConfig(path string) *Config {
	config := new(Config)
	once.Do(func() {
		filePath, err := filepath.Abs(path)
		if err != nil {
			fmt.Print(fmt.Sprintf("[%s]------%c[%d;%d;%dmparse toml error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
			os.Exit(0)
		}
		if _, err := toml.DecodeFile(filePath, config); err != nil {
			fmt.Print(fmt.Sprintf("[%s]------%c[%d;%d;%dmparse toml error:: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
			os.Exit(0)
		}
	})
	return config
}
