package dto

type (
	GenerateParam struct {
		BaseParam
		ProjectName string   `query:"projectName" json:"projectName" validate:"required"` // 项目名
		DataSource  int      `query:"dataSource" json:"dataSource" validate:"required"`   // 项目名
		DbAddr      string   `query:"dbAddr" json:"dbAddr"`                               // 数据库地址
		Username    string   `query:"username" json:"username"`                           // 用户名
		Password    string   `query:"password" json:"password"`                           // 密码
		DbName      string   `query:"dbName" json:"dbName"`                               // 数据库名
		TableNames  []string `query:"tableNames" json:"tableNames"`                       // 表名
	}

	SingleFile struct {
		TemplateName string // 模板名
		FileName     string // 文件名
		Extension    string // 扩展名
	}
)
