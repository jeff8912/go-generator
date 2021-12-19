package handler

import (
	"{{.PackagePrefix}}/constant"
	"{{.PackagePrefix}}/dto"
	"{{.PackagePrefix}}/service"
	"{{.PackagePrefix}}/log"
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
)

type (
	{{ .Struct }}Handler struct {
		CommonHandler
	}
)

var (
	Default{{ .Struct }}Handler = new({{ .Struct }}Handler)
	Page{{ .Struct }}Handler    = new({{ .Struct }}Handler)

	{{ $.LowerCamelStruct}}Service = service.{{ $.Struct }}Service
)

func init() {
	Default{{ .Struct }}Handler.getMapping("{{ .UrlPrefix }}", Default{{ .Struct }}Handler.get)
	Default{{ .Struct }}Handler.postMapping("{{ .UrlPrefix }}", Default{{ .Struct }}Handler.create)
	Default{{ .Struct }}Handler.putMapping("{{ .UrlPrefix }}", Default{{ .Struct }}Handler.update)
	Default{{ .Struct }}Handler.deleteMapping("{{ .UrlPrefix }}", Default{{ .Struct }}Handler.delete)

	Page{{ .Struct }}Handler.getMapping("{{ .UrlPrefix }}/page", Page{{ .Struct }}Handler.page)
}

// 新建
func (p *{{ .Struct }}Handler) create(context echo.Context) error {
    defer p.panicCatch(context)
	res := new(dto.BaseResult)

	body := new(dto.{{ .Struct }})
	errCode := p.readBody(context, body)
	if errCode != constant.SUCCESS {
		return context.JSON(http.StatusOK, res.Padding(errCode))
	}

	res = {{ $.LowerCamelStruct}}Service.Create(body)

	return context.JSON(http.StatusOK, res)
}

// 修改
func (p *{{ .Struct }}Handler) update(context echo.Context) error {
    defer p.panicCatch(context)
	res := new(dto.BaseResult)

	body := new(dto.{{ .Struct }})
	errCode := p.readBody(context, body)
	if errCode != constant.SUCCESS {
		return context.JSON(http.StatusOK, res.Padding(errCode))
	}

	res = {{ $.LowerCamelStruct}}Service.Update(body)

	return context.JSON(http.StatusOK, res)
}

// 查询详情
func (p *{{ .Struct }}Handler) get(context echo.Context) error {
    defer p.panicCatch(context)
	res := new(dto.BaseResult)

	{{ .LowerCamelPk }} {{if eq .PkType "string"}}:= context.QueryParam("{{ .Pk }}"){{ else }}, err := strconv.Atoi(context.QueryParam("{{ .Pk }}")){{end}}
	if {{if eq .PkType "string"}}{{ .LowerCamelPk }} == ""{{ else }}err != nil{{end}} {
		log.Error("[get]参数有误！", "queryString", context.QueryString(), "{{ .LowerCamelPk }}", {{ .LowerCamelPk }}, "err", err)
		return context.JSON(http.StatusOK, res.Padding(constant.BIND_PARAM_ERROR))
	}

	res = {{ $.LowerCamelStruct}}Service.Get({{if eq .PkType "int64"}}int64({{ .LowerCamelPk }}){{ else }}{{ .LowerCamelPk }}{{end}})

	return context.JSON(http.StatusOK, res)
}

// 删除
func (p *{{ .Struct }}Handler) delete(context echo.Context) error {
    defer p.panicCatch(context)
	res := new(dto.BaseResult)

	{{ .LowerCamelPk }} {{if eq .PkType "string"}}:= context.QueryParam("{{ .Pk }}"){{ else }}, err := strconv.Atoi(context.QueryParam("{{ .Pk }}")){{end}}
	if {{if eq .PkType "string"}}{{ .LowerCamelPk }} == ""{{ else }}err != nil{{end}} {
		log.Error("[delete]参数有误！", "queryString", context.QueryString(), "{{ .LowerCamelPk }}", {{ .LowerCamelPk }}, "err", err)
		return context.JSON(http.StatusOK, res.Padding(constant.BIND_PARAM_ERROR))
	}

	res = {{ $.LowerCamelStruct}}Service.Delete({{if eq .PkType "int64"}}int64({{ .LowerCamelPk }}){{ else }}{{ .LowerCamelPk }}{{end}})

	return context.JSON(http.StatusOK, res)
}

// 分页查询
func (p *{{ .Struct }}Handler) page(context echo.Context) error {
    defer p.panicCatch(context)
	res := new(dto.BaseResult)

	body := new(dto.BasePageParam)
	errCode := p.readBody(context, body)
	if errCode != constant.SUCCESS {
		return context.JSON(http.StatusOK, res.Padding(errCode))
	}

	if body.PageSize == 0 {
		body.PageSize = 20
	}
	if body.Page == 0 {
		body.Page = 1
	}

	res = {{ $.LowerCamelStruct}}Service.Page(body)

	return context.JSON(http.StatusOK, res)
}