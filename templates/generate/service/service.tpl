package service

import (
	"{{.PackagePrefix}}/constant"
	"{{.PackagePrefix}}/dto"
	"{{.PackagePrefix}}/modules"

	"{{.PackagePrefix}}/log"
)

type (
	{{ .Struct }}ServiceImpl struct {
	}
)

var (
	{{ .Struct }}Service = new({{ .Struct }}ServiceImpl)
	{{ $.LowerCamelStruct }}Mapper  = modules.{{ .Struct }}Mapper
)

// 新建
func (p *{{ .Struct }}ServiceImpl) Create({{ $.LowerCamelStruct }}Dto *dto.{{ .Struct }}) (res *dto.BaseResult) {
	res = new(dto.BaseResult)

	{{if .HasUniqueIndex}} {{ range $i, $v := .Indexes }}
	{{ $.LowerCamelStruct }}Db, err {{if eq $i 0}}:{{end}}= {{ $.LowerCamelStruct}}Mapper.Get_by{{ range $ii, $vv := $v.Columns }}_{{$vv.ColumnName}}{{end}}({{ range $ii, $vv := $v.Columns }}{{if gt $ii 0}}, {{end}}{{ $.LowerCamelStruct }}Dto.{{$vv.UpperColumn}}{{end}})
	if err != nil {
		log.Error("[Create]根据唯一索引查询失败", "err", err)
		return res.Padding(constant.ADD_ERROR)
	}

	if {{ $.LowerCamelStruct }}Db.{{$.UpperCamelPk}} > 0 {
		res.Code = constant.ADD_ERROR
		res.Msg = "相同{{ range $ii, $vv := $v.Columns }}{{if gt $ii 0}}/{{end}}{{$vv.ColumnComment}}{{end}}已存在!"
		return res
	}
	{{end}}{{end}}

	{{ .LowerCamelStruct }} := modules.{{ .Struct }}{ {{ range $i, $v := .Columns }}{{if eq $i 0}}{{else}}
		{{ $v.UpperColumn }}: {{ $.LowerCamelStruct }}Dto.{{ $v.UpperColumn }},{{end}}{{end}}
	}

	_, err {{if eq .HasUniqueIndex false}}:{{end}}= {{ $.LowerCamelStruct}}Mapper.Create(&{{ .LowerCamelStruct }})
	if err != nil {
		log.Error("[Create]新增失败", "err", err)
		return res.Padding(constant.ADD_ERROR)
	}
	return res.Padding(constant.SUCCESS)
}

// 修改
func (p *{{ .Struct }}ServiceImpl) Update({{ $.LowerCamelStruct }}Dto *dto.{{ .Struct }}) (res *dto.BaseResult) {
	res = new(dto.BaseResult)

	{{if .HasUniqueIndex}} {{ range $i, $v := .Indexes }}
	{{ $.LowerCamelStruct }}Db, err {{if eq $i 0}}:{{end}}= {{ $.LowerCamelStruct}}Mapper.Get_by{{ range $ii, $vv := $v.Columns }}_{{$vv.ColumnName}}{{end}}({{ range $ii, $vv := $v.Columns }}{{if gt $ii 0}}, {{end}}{{ $.LowerCamelStruct }}Dto.{{$vv.UpperColumn}}{{end}})
	if err != nil {
		log.Error("[Update]根据唯一索引查询失败", "err", err)
		return res.Padding(constant.ADD_ERROR)
	}

	if {{ $.LowerCamelStruct }}Db.{{$.UpperCamelPk}} > 0 && {{ $.LowerCamelStruct }}Db.{{$.UpperCamelPk}} != {{ $.LowerCamelStruct }}Dto.{{$.UpperCamelPk}} {
		res.Code = constant.EDIT_ERROR
		res.Msg = "相同{{ range $ii, $vv := $v.Columns }}{{if gt $ii 0}}/{{end}}{{$vv.ColumnComment}}{{end}}已存在!"
		return res
	}
	{{end}}{{end}}

	{{ .LowerCamelStruct }} := modules.{{ .Struct }}{ {{ range $i, $v := .Columns }}{{if eq $i 0}}{{else}}
	{{ $v.UpperColumn }}: {{ $.LowerCamelStruct }}Dto.{{ $v.UpperColumn }},{{end}}{{end}}
	}

	err {{if eq .HasUniqueIndex false}}:{{end}}= {{ $.LowerCamelStruct}}Mapper.Update(&{{ .LowerCamelStruct }})
	if err != nil {
		log.Error("[Update]修改失败", "err", err)
		return res.Padding(constant.ADD_ERROR)
	}
	return res.Padding(constant.SUCCESS)
}

// 查询详情
func (p *{{ .Struct }}ServiceImpl) Get({{ .LowerCamelPk }} {{ .PkType }}) (res *dto.BaseResult) {
	res = new(dto.BaseResult)

	{{ .LowerCamelStruct }}, err := {{ $.LowerCamelStruct}}Mapper.Get({{ .LowerCamelPk }})
	if err != nil {
		log.Error("[Get]查询失败！", "{{ .LowerCamelPk }}", {{ .LowerCamelPk }}, "err", err)
		return res.Padding(constant.QUERY_ERROR)
	}

	res.Data = &dto.{{ .Struct }} { {{ range $i, $v := .Columns }}
		{{ $v.UpperColumn }}: {{ $.LowerCamelStruct }}.{{ $v.UpperColumn }},{{end}}
	}
	return res.Padding(constant.SUCCESS)
}

// 删除
func (p *{{ .Struct }}ServiceImpl) Delete({{ .LowerCamelPk }} {{ .PkType }}) (res *dto.BaseResult) {
	res = new(dto.BaseResult)

	err := {{ $.LowerCamelStruct}}Mapper.Delete({{ .LowerCamelPk }})
	if err != nil {
		log.Error("[delete]删除失败！", "{{ .LowerCamelPk }}", {{ .LowerCamelPk }}, "err", err)
		return res.Padding(constant.DELETE_ERROR)
	}
	return res.Padding(constant.SUCCESS)
}

// 分页查询
func (p *{{ .Struct }}ServiceImpl) Page(pageParam *dto.BasePageParam) (res *dto.BaseResult) {
    res = new(dto.BaseResult)
    pageResult := new(dto.PageResult)

    queryParams := map[string]interface{}{}
    queryParams["page"] = pageParam.Page
    queryParams["pageSize"] = pageParam.PageSize

    records, total, err := {{ $.LowerCamelStruct}}Mapper.Page(queryParams)
    if err != nil {
        log.Error("[Page]分页失败！", "pageParam", pageParam, "err", err)
        return res.Padding(constant.PAGE_ERROR)
    }

    pageResult.Items = records
    pageResult.Total = total
    res.Data = pageResult
    return res.Padding(constant.SUCCESS)
}