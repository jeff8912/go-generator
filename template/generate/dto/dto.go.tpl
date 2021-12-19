package dto

{{if .HasTime}}import "time"{{end}}

type (
	{{ .Struct }} struct { {{ range $i, $v := .Columns }}
		{{ $v.UpperColumn }} {{ $v.ColumnType }} `json:"{{ $v.JsonColumn }}"` // {{ $v.ColumnComment }} {{end}}
	}
)
