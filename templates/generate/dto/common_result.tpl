package dto

import "{{.packagePrefix}}/constant"

type (
	//基础响应
	BaseResult struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	//分页结果
	PageResult struct {
		Total int         `json:"total"`
		Items interface{} `json:"items"`
	}
)

var (
	ResultMap map[int]string
)

func init() {
	ResultMap = make(map[int]string)
	ResultMap[constant.SUCCESS] = "成功"
	ResultMap[constant.PANIC_ERROR] = "PANIC异常"
	ResultMap[constant.BIND_PARAM_ERROR] = "绑定参数异常"
	ResultMap[constant.VALIDATE_PARAM_ERROR] = "校验参数异常"
	ResultMap[constant.READ_PARAM_ERROR] = "读取参数异常"
	ResultMap[constant.UNMARSHAL_PARAM_ERROR] = "反序列化参数异常"
	ResultMap[constant.MARSHAL_PARAM_ERROR] = "序列化参数异常"
	ResultMap[constant.ENCRYPT_ERROR] = "加密错误"
	ResultMap[constant.DECRYPT_ERROR] = "解密错误"
	ResultMap[constant.SIGN_FAIL] = "签名校验失败"
	ResultMap[constant.SERVER_ERROR] = "服务器异常"
	ResultMap[constant.DB_ERROR] = "数据库异常"
	ResultMap[constant.REDIS_ERROR] = "redis缓存异常"
	ResultMap[constant.UNKNOWN] = "未知错误"
	ResultMap[constant.JSON_MARSHAL_ERROR] = "JSON序列化失败"
	ResultMap[constant.JSON_UNMARSHAL_ERROR] = "JSON反序列化失败"
	ResultMap[constant.ADD_ERROR] = "新增失败"
	ResultMap[constant.DELETE_ERROR] = "删除失败"
	ResultMap[constant.EDIT_ERROR] = "修改失败"
	ResultMap[constant.QUERY_ERROR] = "查询失败"
	ResultMap[constant.PAGE_ERROR] = "分页查询失败"
}

func (res *BaseResult) Padding(code int) *BaseResult {
	res.Code = code
	if val, ok := ResultMap[code]; ok {
		res.Msg = val
	} else {
		val, _ := ResultMap[constant.UNKNOWN]
		res.Msg = val
	}
	return res
}

func (res *BaseResult) Padding2(code int, msg string) *BaseResult {
	res.Code = code
	res.Msg = msg
	return res
}

func (res *BaseResult) Padding3(code int, msg string, data interface{}) *BaseResult {
	res.Code = code
	res.Msg = msg
	res.Data = data
	return res
}

func (res *BaseResult) PaddingData(data interface{}) *BaseResult {
	res.Data = data
	return res
}
