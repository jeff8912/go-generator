package handlers

import (
	"{{.PackagePrefix}}/constant"
	"{{.PackagePrefix}}/dto"
	"{{.PackagePrefix}}/service"
	"net/http"
	"{{.PackagePrefix}}/log"
	"strconv"

	"github.com/labstack/echo"

	"runtime/debug"
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
	res := new(dto.BaseResult)

	defer func() {
		if err := recover(); err != nil {
			log.Error("[create]panic", "queryString", context.QueryString(),
				"err", err, "stack", string(debug.Stack()))
			context.JSON(http.StatusOK, res.Padding(constant.PANIC_ERROR))
		}
	}()

	body := new(dto.{{ .Struct }})
	errCode := p.readBody(context, body)
	if errCode != constant.SUCCESS {
		return context.JSON(http.StatusOK, res.Padding(errCode))
	}

	log.Info("[create]begin", "queryString", context.QueryString(), "body", body)

	res = {{ $.LowerCamelStruct}}Service.Create(body)

	if log.GetLogLevel() == "debug" {
		log.Debug("[create]end", "queryString", context.QueryString(), "response", res)
	}
	return context.JSON(http.StatusOK, res)
}

// 修改
func (p *{{ .Struct }}Handler) update(context echo.Context) error {
	res := new(dto.BaseResult)

	defer func() {
		if err := recover(); err != nil {
			log.Error("[update]panic", "queryString", context.QueryString(),
				"err", err, "stack", string(debug.Stack()))
			context.JSON(http.StatusOK, res.Padding(constant.PANIC_ERROR))
		}
	}()

	body := new(dto.{{ .Struct }})
	errCode := p.readBody(context, body)
	if errCode != constant.SUCCESS {
		return context.JSON(http.StatusOK, res.Padding(errCode))
	}

	log.Info("[update]begin", "queryString", context.QueryString(), "body", body)

	res = {{ $.LowerCamelStruct}}Service.Update(body)

	if log.GetLogLevel() == "debug" {
		log.Debug("[update]end", "queryString", context.QueryString(), "response", res)
	}
	return context.JSON(http.StatusOK, res)
}

// 查询详情
func (p *{{ .Struct }}Handler) get(context echo.Context) error {
	res := new(dto.BaseResult)

	defer func() {
		if err := recover(); err != nil {
			log.Error("[get]panic", "queryString", context.QueryString(),
			"err", err, "stack", string(debug.Stack()))
			context.JSON(http.StatusOK, res.Padding(constant.PANIC_ERROR))
		}
	}()

	log.Info("[get]begin", "queryString", context.QueryString())

	{{ .LowerCamelPk }} {{if eq .PkType "string"}}:= context.QueryParam("{{ .Pk }}"){{ else }}, err := strconv.Atoi(context.QueryParam("{{ .Pk }}")){{end}}
	if {{if eq .PkType "string"}}{{ .LowerCamelPk }} == ""{{ else }}err != nil{{end}} {
		log.Error("[get]参数有误！", "queryString", context.QueryString(), "{{ .LowerCamelPk }}", {{ .LowerCamelPk }}, "err", err)
		return context.JSON(http.StatusOK, res.Padding(constant.BIND_PARAM_ERROR))
	}

	res = {{ $.LowerCamelStruct}}Service.Get({{if eq .PkType "int64"}}int64({{ .LowerCamelPk }}){{ else }}{{ .LowerCamelPk }}{{end}})

	if log.GetLogLevel() == "debug" {
		log.Debug("[get]end", "queryString", context.QueryString(), "response", res)
	}
	return context.JSON(http.StatusOK, res)
}

// 删除
func (p *{{ .Struct }}Handler) delete(context echo.Context) error {
	res := new(dto.BaseResult)

	defer func() {
		if err := recover(); err != nil {
			log.Error("[delete]panic", "queryString", context.QueryString(),
				"err", err, "stack", string(debug.Stack()))
			context.JSON(http.StatusOK, res.Padding(constant.PANIC_ERROR))
		}
	}()

	log.Info("[delete]begin", "queryString", context.QueryString())

	{{ .LowerCamelPk }} {{if eq .PkType "string"}}:= context.QueryParam("{{ .Pk }}"){{ else }}, err := strconv.Atoi(context.QueryParam("{{ .Pk }}")){{end}}
	if {{if eq .PkType "string"}}{{ .LowerCamelPk }} == ""{{ else }}err != nil{{end}} {
		log.Error("[delete]参数有误！", "queryString", context.QueryString(), "{{ .LowerCamelPk }}", {{ .LowerCamelPk }}, "err", err)
		return context.JSON(http.StatusOK, res.Padding(constant.BIND_PARAM_ERROR))
	}

	res = {{ $.LowerCamelStruct}}Service.Delete({{if eq .PkType "int64"}}int64({{ .LowerCamelPk }}){{ else }}{{ .LowerCamelPk }}{{end}})

	if log.GetLogLevel() == "debug" {
		log.Debug("[delete]end", "queryString", context.QueryString(), "response", res)
	}
	return context.JSON(http.StatusOK, res)
}

// 分页查询
func (p *{{ .Struct }}Handler) page(context echo.Context) error {
	res := new(dto.BaseResult)

	defer func() {
		if err := recover(); err != nil {
			log.Error("[page]panic", "queryString", context.QueryString(),
				"err", err, "stack", string(debug.Stack()))
			context.JSON(http.StatusOK, res.Padding(constant.PANIC_ERROR))
		}
	}()

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

	log.Info("[page]begin", "queryString", context.QueryString(), "body", body)

	res = {{ $.LowerCamelStruct}}Service.Page(body)

	if log.GetLogLevel() == "debug" {
		log.Debug("[page]end", "queryString", context.QueryString(), "response", res)
	}
	return context.JSON(http.StatusOK, res)
}