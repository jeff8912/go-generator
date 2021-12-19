package handler

import (
	"github.com/labstack/echo/v4"
	"go-generator/constant"
	"go-generator/dto"
	"go-generator/service"
)

type (
	DbHandlerImpl struct {
		CommonHandler
	}
)

var (
	dbService       = service.DbService
	DatabaseHandler = new(DbHandlerImpl)
	TableHandler    = new(DbHandlerImpl)
)

func init() {
	DatabaseHandler.postMapping("databases", DatabaseHandler.databases)
	TableHandler.postMapping("tables", TableHandler.tables)
}

func (p *DbHandlerImpl) databases(context echo.Context) error {
	defer p.panicCatch(context)
	var (
		res     = new(dto.BaseResult)
		body    = new(dto.DbQueryParam)
		errCode int
	)
	errCode = p.readBody(context, body)
	if errCode != constant.SUCCESS {
		return p.writeBody(context, res.Padding(errCode))
	}

	res = dbService.GetDatabases(body)
	return p.writeBody(context, res)
}

func (p *DbHandlerImpl) tables(context echo.Context) error {
	defer p.panicCatch(context)
	var (
		res     = new(dto.BaseResult)
		body    = new(dto.TbQueryParam)
		errCode int
	)
	errCode = p.readBody(context, body)
	if errCode != constant.SUCCESS {
		return p.writeBody(context, res.Padding(errCode))
	}

	res = dbService.GetTables(body)
	return p.writeBody(context, res)
}
