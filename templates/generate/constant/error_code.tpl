package constant

const (
	// http response status_code
	SUCCESS               = 0  // 成功
	PANIC_ERROR           = 1  // PANIC异常
	BIND_PARAM_ERROR      = 2  // 绑定参数异常
	VALIDATE_PARAM_ERROR  = 3  // 校验参数异常
	READ_PARAM_ERROR      = 4  // 读取参数异常
	UNMARSHAL_PARAM_ERROR = 5  // 反序列化参数异常
	MARSHAL_PARAM_ERROR   = 6  // 序列化参数异常
	ENCRYPT_ERROR         = 7  // 加密错误
	DECRYPT_ERROR         = 8  // 解密错误
	SIGN_FAIL             = 9  // 签名校验失败
	SERVER_ERROR          = 10 // 服务器异常
	DB_ERROR              = 11 // 数据库异常
	REDIS_ERROR           = 12 // redis缓存异常
	UNKNOWN               = 13 // 未知错误
	JSON_MARSHAL_ERROR    = 14 // JSON序列化失败
	JSON_UNMARSHAL_ERROR  = 15 // JSON反序列化失败
    ADD_ERROR             = 17 // 新增失败
    DELETE_ERROR          = 18 // 删除失败
    EDIT_ERROR            = 19 // 修改失败
    QUERY_ERROR           = 20 // 查询失败
    PAGE_ERROR            = 21 // 分页查询失败
)
