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
		Total int64       `json:"total"`
		Items interface{} `json:"items"`
	}
)

func (res *BaseResult) Padding(code int) *BaseResult {
	res.Code = code
	if val, ok := constant.ErrCodeMap[code]; ok {
		res.Msg = val
	} else {
		val, _ := constant.ErrCodeMap[constant.UNKNOWN]
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
