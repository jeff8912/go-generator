package handler

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"go-generator/constant"
	"go-generator/dto"
	"go-generator/log"
	"go-generator/service"
	"go-generator/util"
	"net/http"
	"strings"
)

type (
	GenerateHandlerImpl struct {
		CommonHandler
	}
)

var (
	generateService = service.GenerateService
	GenerateHandler = new(GenerateHandlerImpl)
)

func init() {
	GenerateHandler.postMapping("generate", GenerateHandler.generate)
}

func (p *GenerateHandlerImpl) generate(context echo.Context) error {
	defer p.panicCatch(context)

	var (
		res     = new(dto.BaseResult)
		body    = new(dto.GenerateParam)
		errCode int
	)
	errCode = p.readBody(context, body)
	if errCode != constant.SUCCESS {
		return p.writeBody(context, res.Padding(errCode))
	}

	projectNames := strings.Split(body.ProjectName, "/")
	zipPath := projectNames[0]
	if zipPath == "etc" || zipPath == "public" || zipPath == "template" {
		log.Error("[代码生成]项目名不能为etc/public/template")
		return p.writeBody(context, res.Padding(constant.VALIDATE_PARAM_ERROR))
	}

	if body.DataSource == constant.DATASOURCE_MYSQL && (body.DbAddr == "" || body.Username == "" || body.Password == "") {
		log.Error("请先确认数据库配置")
		return p.writeBody(context, res.Padding(constant.VALIDATE_PARAM_ERROR))
	}

	zipFile := projectNames[0] + ".zip"

	err := generateService.Generate(body)
	if err == nil {
		buffer, err := util.CreateZip(zipPath, zipFile)
		if err != nil {
			log.Error("[代码生成]创建压缩包失败", "err", err)
			return p.writeBody(context, res.Padding(constant.SERVER_ERROR))
		}

		util.RemovePath(zipPath, zipFile)

		//设置请求头  使用浏览器下载
		context.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+zipFile)
		return context.Stream(http.StatusOK, echo.MIMEOctetStream, bytes.NewReader(buffer))
	} else {
		log.Error("[代码生成]失败", "err", err)
	}

	util.RemovePath(zipPath, zipFile)

	return p.writeBody(context, res.Padding(constant.SERVER_ERROR))
}
