package dto

type (
	DbQueryParam struct {
		BaseParam
		DbAddr   string `query:"dbAddr" json:"dbAddr" validate:"required"`     // 数据库地址
		Username string `query:"username" json:"username" validate:"required"` // 用户名
		Password string `query:"password" json:"password" validate:"required"` // 密码
	}

	Database struct {
		Id          int    `json:"id"`
		TableSchema string `json:"name" gorm:"column:table_schema"` // 数据库名称
	}
)
