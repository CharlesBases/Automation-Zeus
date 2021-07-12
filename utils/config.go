package utils

import (
	"os"
	"path/filepath"
	"sync"

	"mysql-gen-go/logger"

	"github.com/BurntSushi/toml"
)

var once sync.Once

// ParseConfig 解析配置文件
func ParseConfig(path string) *Config {
	config := new(Config)
	once.Do(func() {
		filePath, err := filepath.Abs(path)
		if err != nil {
			logger.Error("parse toml error: %v", err)
			os.Exit(0)
		}
		if _, err := toml.DecodeFile(filePath, config); err != nil {
			logger.Error("parse toml error: %v", err)
			os.Exit(0)
		}
	})
	return config
}
