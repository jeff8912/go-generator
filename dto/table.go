package dto

type (
	TbQueryParam struct {
		BaseParam
		DbAddr   string `query:"dbAddr" json:"dbAddr" validate:"required"`     // 数据库地址
		Username string `query:"username" json:"username" validate:"required"` // 用户名
		Password string `query:"password" json:"password" validate:"required"` // 密码
		DbName   string `query:"dbName" json:"dbName" validate:"required"`     // 数据库名
	}

	TableOptions struct {
		TableName    string `json:"name" gorm:"column:table_name"`       // 表名
		TableComment string `json:"comment" gorm:"column:table_comment"` // 表注释
	}

	Column struct {
		TableName     string `json:"table_name" gorm:"column:table_name"`         // 表名
		ColumnName    string `json:"column_name" gorm:"column:column_name"`       // 字段名
		DataType      string `json:"data_type" gorm:"column:data_type"`           // 字段类型
		ColumnComment string `json:"column_comment" gorm:"column:column_comment"` // 字段注释
		IsNullable    string `json:"is_nullable" gorm:"column:is_nullable"`       // 列的为空性。如果列允许 NULL，那么该列返回 YES。否则，返回 NO。
		UpperColumn   string // 大驼峰字段
		LowerColumn   string // 小驼峰字段
		JsonColumn    string // json字段
		ColumnType    string // 字段类型
	}
	TableColumn struct {
		Columns          []*Column `json:"columns"`          // 表字段
		NoPkColumns      []*Column `json:"noPkColumns"`      // 不带主键的表字段
		Struct           string    `json:"struct"`           // 表结构体定义
		LowerCamelStruct string    `json:"lowerCamelStruct"` // 驼峰结构体
		UrlPrefix        string    `json:"urlPrefix"`        // restful路径前缀
		TableName        string    `json:"tableName"`        // 表名
		TableComment     string    `json:"tableComment"`     // 表注释
		Index            int       `json:"index"`            // 序号
		Pk               string    `json:"pk"`               // 主键
		PkType           string    `json:"pkType"`           // 主键类型
		LowerCamelPk     string    `json:"lowerCamelPk"`     // 驼峰主键
		UpperCamelPk     string    `json:"upperCamelPk"`     // 大驼峰主键
		HasTime          bool      `json:"hasTime"`          // 是否包含时间字段
		PackagePrefix    string    `json:"packagePrefix"`    // 包前缀
		HasUniqueIndex   bool      `json:"hasUniqueIndex"`   // 是否包含唯一索引
		Indexes          []Index
	}
	TableIndex struct {
		Table      string `json:"table" gorm:"column:Table"`               // 表名
		NonUnique  int    `json:"non_unique" gorm:"column:Non_unique"`     // 非唯一键
		KeyName    string `json:"key_name" gorm:"column:Key_name"`         // 索引名称
		SeqInIndex string `json:"seq_in_index" gorm:"column:Seq_in_index"` // 在当前索引中的序号
		ColumnName string `json:"column_name" gorm:"column:Column_name"`   // 字段名称
	}
	Index struct {
		KeyName string // 索引名称
		Columns []Column
	}
)
