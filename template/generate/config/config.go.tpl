package config

import (
	"github.com/Unknwon/goconfig"
	"log"
	"os"
	"path/filepath"
)

var defaultCfg *goconfig.ConfigFile

func init() {
	fileAbsPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}

	appPath := filepath.Dir(fileAbsPath)
	// 支持go run 和 go test
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	workPath, err = filepath.Abs(workPath)
	if err != nil {
		panic(err)
	}

	configPath := "etc/conf.ini"

	// 现在应用对应目录查找
	cfg, err := goconfig.LoadConfigFile(appPath + "/" + configPath)

	// 查找失败，在当前目录查找
	if err != nil {
		cfg, err = goconfig.LoadConfigFile(workPath + "/" + configPath)
	}

	// 当前目录查找失败，在去当前目录上一级目录查找
	if err != nil {
		cfg, err = goconfig.LoadConfigFile(filepath.Dir(workPath) + "/" + configPath)
	}

	// 都没有找到，返回失败
	if err != nil {
		log.Fatalf("load config conf.ini failed, error = %s!", err)
	}

	defaultCfg = cfg
}

func GetValue(section, key string) string {
	value, err := defaultCfg.GetValue(section, key)
	if err != nil {
		log.Fatalf("section = %s, get item key = %s failed, err = %s!", section, key, err)
	}
	return value
}

func GetSections() []string {
	return defaultCfg.GetSectionList()
}

func GetDefaultConfig() *goconfig.ConfigFile {
	return defaultCfg
}

func Int(section, key string) int {
	value, err := defaultCfg.Int(section, key)
	if err != nil {
		log.Fatalf("section = %s, get item key = %s failed, Int val err = %s!", section, key, err)
	}
	return value
}
