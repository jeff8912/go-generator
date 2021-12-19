package module

import (
	"errors"
	{{if .HasUniqueIndex}}"gorm.io/gorm"{{end}}
	"{{.PackagePrefix}}/log"
	{{if .HasTime}}"time"{{end}}
)

// {{ .TableName }} {{ .TableComment }}
type {{ .Struct }} struct { {{ range $i, $v := .Columns }}
	{{ $v.UpperColumn }} {{ $v.ColumnType }} {{if eq $i 0}} `gorm:"column:{{ $v.ColumnName }};not null;primary_key;AUTO_INCREMENT" json:"{{ $v.JsonColumn }}"` {{else}} `gorm:"column:{{ $v.ColumnName }}" json:"{{ $v.JsonColumn }}"` {{end}} // {{ $v.ColumnComment }} {{end}}
}

var (
	{{ .Struct }}Mapper = new({{ .Struct }})
)

func (*{{ .Struct }}) TableName() string {
	return "{{ .TableName }}"
}

// 根据主键查询记录
func (p *{{ .Struct }}) Get({{ .LowerCamelPk }} {{ .PkType }}) (*{{ .Struct }}, error) {
	record := new({{ .Struct }})
	db := DbInstance.Debug()
	res := db.Where("{{ .Pk }} = ?", {{ .LowerCamelPk }}).First(record)
	return record, res.Error
}

{{if .HasUniqueIndex}} {{ range $i, $v := .Indexes }}
func (p *{{ $.Struct }}) GetBy{{ range $ii, $vv := $v.Columns }}{{$vv.UpperColumn}}{{end}}({{ range $ii, $vv := $v.Columns }}{{if gt $ii 0}}, {{end}}{{$vv.LowerColumn}} {{$vv.ColumnType}}{{end}}) (*{{ $.Struct }}, error) {
	record := new({{ $.Struct }})
	db := DbInstance.Debug()
	res := db.Where("{{ range $ii, $vv := $v.Columns }}{{if gt $ii 0}} AND {{end}}{{$vv.ColumnName}} = ?{{end}}"{{ range $ii, $vv := $v.Columns }}, {{$vv.LowerColumn}}{{end}}).First(record)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return record, res.Error
	}
	return record, nil
}
{{end}}{{end}}

// 新增记录
func (p *{{ .Struct }}) Create({{ .LowerCamelStruct }} *{{ .Struct }}) (*{{ .Struct }}, error) {
	db := DbInstance.Debug()
	res := db.Create({{ .LowerCamelStruct }})
	if res.Error != nil {
		log.Error("[{{ .TableName }}_module][创建失败]", "err", res.Error)
		return nil, errors.New("创建失败")
	}
	return {{ .LowerCamelStruct }}, nil
}

// 修改记录
func (p *{{ .Struct }}) Update({{ .LowerCamelStruct }} *{{ .Struct }}) error {
	param := make(map[string]interface{})
	{{range $i, $v := .Columns}}{{if eq $v.ColumnName "create_time"}}{{else if eq $v.ColumnName "create_date"}}{{else if eq $v.ColumnName "create_by"}}{{else}}
	param["{{ $v.ColumnName }}"] = {{ $.LowerCamelStruct }}.{{ $v.UpperColumn }}{{end}}{{end}}
	res := DbInstance.Debug().Table(p.TableName()).
		Where("{{ .Pk }} = ?", param["{{ .Pk }}"]).
		Updates(param)
	if res.Error != nil && res.Error != UPDATE_AFFECTED_ZERO_ERROR {
		log.Error("[{{ .TableName }}_module][修改失败]", "err", res.Error)
		return errors.New("修改失败")
	}
	return nil
}

// 删除记录
func (p *{{ .Struct }}) Delete({{ .LowerCamelPk }} {{ .PkType }}) error {
	db := DbInstance.Debug()
	res := db.Where("{{ .Pk }} = ?", {{ .LowerCamelPk }}).
		Delete(&{{ .Struct }}{})
	return res.Error
}

// 查询所有记录
func (p *{{ .Struct }}) List() ([]{{ .Struct }}, error) {
	records := []{{ .Struct }}{}
	db := DbInstance.Debug()
	res := db.Find(&records)
	return records, res.Error
}

// 分页查询记录
func (p *{{ .Struct }}) Page(queryParams map[string]interface{}) ([]{{ .Struct }}, int64, error) {
	records := []{{ .Struct }}{}
	db := DbInstance.Debug()
	queryTotal := db.Model(&{{ .Struct }}{})
	db = setPage(queryParams, db)
	total := getTotal(queryParams, queryTotal)

	res := db.Where(queryParams).Find(&records)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
        return records, total, res.Error
    }
	return records, total, nil
}