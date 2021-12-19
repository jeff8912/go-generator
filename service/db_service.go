package service

import (
	"go-generator/constant"
	"go-generator/dto"
	"go-generator/log"
	"go-generator/module"
)

type (
	DbServiceImpl struct {
		CommonServiceImpl
	}
)

var (
	DbService = new(DbServiceImpl)
)

func (p *DbServiceImpl) GetDatabases(param *dto.DbQueryParam) *dto.BaseResult {
	res := new(dto.BaseResult)

	db, err := module.GetDbInstance(param.DbAddr, param.Username, param.Password, "")
	if err != nil {
		log.Error("[代码生成]获取数据库连接失败", "err", err)
		return res.Padding(constant.DB_ERROR)
	}

	databases, err := module.GetDatabases(db)
	if err != nil {
		log.Error("[代码生成]查询数据库失败", "err", err)
		return res.Padding(constant.DB_ERROR)
	}
	res.Data = databases

	return res
}

func (p *DbServiceImpl) GetTables(param *dto.TbQueryParam) *dto.BaseResult {
	res := new(dto.BaseResult)

	db, err := module.GetDbInstance(param.DbAddr, param.Username, param.Password, "")
	if err != nil {
		log.Error("[代码生成]获取数据库连接失败", "err", err)
		return res.Padding(constant.DB_ERROR)
	}

	tables, err := module.GetTables(db, param.DbName)
	if err != nil {
		log.Error("[代码生成]查询数据库表失败", "err", err)
		return res.Padding(constant.DB_ERROR)
	}
	res.Data = tables

	return res
}
