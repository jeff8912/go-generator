package dto

type (
	//基础请求实体
	BaseParam struct {
	}

	//基础分页请求实体
	BasePageParam struct {
		Page     int `query:"page" json:"page"`
		PageSize int `query:"pageSize" json:"pageSize"`
	}
)
