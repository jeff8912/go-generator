package service

import (
	"go-generator/constant"
	"go-generator/dto"
	"testing"
)

func TestGenerateFromTemplate(t *testing.T) {
	// 去掉盘符的绝对模板路径
	templateRootPath = "/workspace/go-generator/template/generate/"
	err := generateFromTemplate("ad/test", "etc/conf.ini.tpl", "", nil)
	if err != nil {
		t.Log("err", err)
	}
}

func TestGenerate(t *testing.T) {
	templateRootPath = "/workspace/go-generator/template/generate/"
	initTplFiles()
	param := new(dto.GenerateParam)
	param.ProjectName = "test"
	param.DataSource = constant.DATASOURCE_MYSQL
	param.DbAddr = "127.0.0.1:3306"
	param.DbName = "test"
	param.Username = "root"
	param.Password = "root"
	GenerateService.Generate(param)
}
