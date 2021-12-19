{{range $i, $v := .tableColumns }}# {{$v.Index}} {{$v.TableComment}}
## {{$v.Index}}.1 创建

- 请求路径：

```
{{$v.UrlPrefix}}
```

- 请求方式：

POST

- 请求参数：

| 字段名   | 字段类型               | 必填 | 说明                   | 示例                         |
| -------- | ---------- | ---- | ---------------------- | ------------------- |{{range $index, $value := .NoPkColumns }}{{if eq $value.ColumnName "create_time"}}{{else if eq $value.ColumnName "create_date"}}{{else if eq $value.ColumnName "create_by"}}{{else if eq $value.ColumnName "update_time"}}{{else if eq $value.ColumnName "update_time"}}{{else if eq $value.ColumnName "update_time"}}{{else}}
| {{$value.JsonColumn}}| {{$value.ColumnType}} | {{if eq $value.IsNullable "NO"}}Y{{else}}N{{end}}|{{$value.ColumnComment}}||{{end}}{{end}}

- 响应参数

| 字段名 | 字段类型             | 说明                 | 示例        |
| ------ | ------------------| --------------------| ----------- |
| code   | int               | 成功:0,失败：非 0     | 0           |
| data   | Object            | 返回的数据            |             |
| msg    | string            | 返回错误信息          | 失败         |

## {{$v.Index}}.2 修改

- 请求路径：

```
{{$v.UrlPrefix}}
```

- 请求方式：

PUT

- 请求参数：

| 字段名   | 字段类型               | 必填 | 说明                   | 示例                         |
| -------- | ---------- | ---- | ---------------------- | ------------------- |{{range $index, $value := .Columns }}{{if eq $value.ColumnName "create_time"}}{{else if eq $value.ColumnName "create_date"}}{{else if eq $value.ColumnName "create_by"}}{{else if eq $value.ColumnName "update_time"}}{{else if eq $value.ColumnName "update_time"}}{{else if eq $value.ColumnName "update_time"}}{{else}}
| {{$value.JsonColumn}}| {{$value.ColumnType}} | {{if eq $value.IsNullable "NO"}}Y{{else}}N{{end}}|{{$value.ColumnComment}}||{{end}}{{end}}

- 响应参数

| 字段名 | 字段类型             | 说明                 | 示例        |
| ------ | ------------------| --------------------| ----------- |
| code   | int               | 成功:0,失败：非 0     | 0           |
| data   | Object            | 返回的数据            |             |
| msg    | string            | 返回错误信息          | 失败         |

## {{$v.Index}}.3 详情

- 请求路径：

```
{{$v.UrlPrefix}}
```

- 请求方式：

GET

- 请求参数：

| 字段名   | 字段类型               | 必填 | 说明                   | 示例                         |
| -------- | ---------- | ---- | ---------------------- | ------------------- |
{{range $index, $value := .Columns }}{{if eq $index 0}}| {{$value.JsonColumn}}| {{$value.ColumnType}} | Y|{{$value.ColumnComment}}||{{end}}{{end}}

- 响应参数

| 字段名 | 字段类型             | 说明                 | 示例        |
| ------ | ------------------| --------------------| ----------- |
| code   | int               | 成功:0,失败：非 0     | 0           |
| data   | Object{{$.lt}}{{$v.Struct}}>            | 返回的数据            |             |
| msg    | string            | 返回错误信息          | 失败         |

- {{$v.Struct}}

| 字段名   | 字段类型               |  说明                   | 示例                         |
| -------- | -----------------   |  ---------------------- | ------------------- |{{range $index, $value := .Columns }}
| {{$value.JsonColumn}}| {{$value.ColumnType}} | {{$value.ColumnComment}}||{{end}}

## {{$v.Index}}.4 删除

- 请求路径：

```
{{$v.UrlPrefix}}
```

- 请求方式：

DELETE

- 请求参数：

| 字段名   | 字段类型               | 必填 | 说明                   | 示例                         |
| -------- | ---------- | ---- | ---------------------- | ------------------- |
{{range $index, $value := .Columns }}{{if eq $index 0}}| {{$value.JsonColumn}}| {{$value.ColumnType}} | Y|{{$value.ColumnComment}}||{{end}}{{end}}

- 响应参数

| 字段名 | 字段类型             | 说明                 | 示例        |
| ------ | ------------------| --------------------| ----------- |
| code   | int               | 成功:0,失败：非 0     | 0           |
| data   | Object            | 返回的数据            |             |
| msg    | string            | 返回错误信息          | 失败         |

## {{$v.Index}}.5 分页查询

- 请求路径：

```
{{$v.UrlPrefix}}/page
```

- 请求方式：

GET

- 请求参数：

| 字段名    | 字段类型     | 必填 | 说明                   | 示例       |
| ---------| ---------- | ----| ---------------------- | ----------|
| page     |    int     |   Y | 当前页                  |   1       |
| pageSize     |    int     |   Y | 每页显示数量             |   10      |

- 响应参数

| 字段名 | 字段类型             | 说明                 | 示例        |
| ------ | ------------------| --------------------| ----------- |
| code   | int               | 成功:0,失败：非 0     | 0           |
| data   | Object{{$.lt}}{{$v.Struct}}>            | 返回的数据            |             |
| msg    | string            | 返回错误信息          | 失败         |

- {{$v.Struct}}

| 字段名   | 字段类型               |  说明                   | 示例                         |
| -------- | -----------------   |  ---------------------- | ------------------- |{{range $index, $value := .Columns }}
| {{$value.JsonColumn}}| {{$value.ColumnType}} | {{$value.ColumnComment}}||{{end}}
{{end}}
