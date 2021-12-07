package config

import (
	"{{.packagePrefix}}/common"
	"log"
	"path/filepath"

	"github.com/Unknwon/goconfig"
)

var defaultCfg *goconfig.ConfigFile

func init() {
	appPath := common.GetMainPath()
	// 支持go run 和 go test
	workPath := common.GetWorkPath()
	configPath := "etc/conf.ini"

	// 现在应用对应目录查找
	cfg, err := goconfig.LoadConfigFile(common.JoinPath(appPath, configPath))

	// 查找失败，在当前目录查找
	if err != nil {
		cfg, err = goconfig.LoadConfigFile(common.JoinPath(workPath, configPath))
	}

	// 当前目录查找失败，在去当前目录上一级目录查找
	if err != nil {
		cfg, err = goconfig.LoadConfigFile(common.JoinPath(filepath.Dir(workPath), configPath))
	}

	// 都没有找到，返回失败
	if err != nil {
		programName := common.GetMainName()
		log.Fatalf("program %s, load config conf.ini failed, error = %s!", programName, err)
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

func GetDefaultCofig() *goconfig.ConfigFile {
	return defaultCfg
}

func Int(section, key string) int {
	value, err := defaultCfg.Int(section, key)
	if err != nil {
		log.Fatalf("section = %s, get item key = %s failed, Int val err = %s!", section, key, err)
	}
	return value
}
