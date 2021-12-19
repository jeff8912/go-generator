package constant

const (
	SUCCESS               = 0   //成功
	SERVER_ERROR          = 1   //服务器错误
	VALIDATE_ERROR        = 2   //校验参数失败
	ENCRYPT_ERROR         = 3   //加密失败
	DECRYPT_ERROR         = 4   //解密失败
	SIGN_ERROR            = 5   //签名失败
	DB_ERROR              = 6   //数据库错误
	REDIS_ERROR           = 7   //缓存失败
	TOKEN_ERROR           = 8   //token无效
	TOKEN_NULL            = 9   //token不能为空
	PANIC_ERROR           = 10  // PANIC异常
	BIND_PARAM_ERROR      = 11  // 绑定参数异常
	VALIDATE_PARAM_ERROR  = 12  // 校验参数异常
	READ_PARAM_ERROR      = 13  // 读取参数异常
	UNMARSHAL_PARAM_ERROR = 14  // 反序列化参数异常
	MARSHAL_PARAM_ERROR   = 15  // 序列化参数异常
	SIGN_FAIL             = 16  // 签名校验失败
	UNKNOWN               = 17  // 未知错误
	JSON_MARSHAL_ERROR    = 18  // JSON序列化失败
	JSON_UNMARSHAL_ERROR  = 19  // JSON反序列化失败
	ADD_ERROR             = 20  // 新增失败
	DELETE_ERROR          = 21  // 删除失败
	EDIT_ERROR            = 22  // 修改失败
	QUERY_ERROR           = 23  // 查询失败
	PAGE_ERROR            = 24  // 分页查询失败
	JSON_ERROR            = 101 //json参数格式有误
)

var (
	ErrCodeMap = map[int]string{
		SUCCESS:               "成功",
		SERVER_ERROR:          "服务器错误",
		VALIDATE_ERROR:        "校验参数失败",
		ENCRYPT_ERROR:         "加密失败",
		DECRYPT_ERROR:         "解密失败",
		SIGN_ERROR:            "签名错误",
		DB_ERROR:              "数据库错误",
		REDIS_ERROR:           "缓存失败",
		TOKEN_ERROR:           "登录已过期",
		TOKEN_NULL:            "token不能为空",
		PANIC_ERROR:           "PANIC异常",
		BIND_PARAM_ERROR:      "绑定参数异常",
		VALIDATE_PARAM_ERROR:  "校验参数异常",
		READ_PARAM_ERROR:      "读取参数异常",
		UNMARSHAL_PARAM_ERROR: "反序列化参数异常",
		MARSHAL_PARAM_ERROR:   "序列化参数异常",
		SIGN_FAIL:             "签名校验失败",
		UNKNOWN:               "未知错误",
		JSON_MARSHAL_ERROR:    "JSON序列化失败",
		JSON_UNMARSHAL_ERROR:  "JSON反序列化失败",
		ADD_ERROR:             "新增失败",
		DELETE_ERROR:          "删除失败",
		EDIT_ERROR:            "修改失败",
		QUERY_ERROR:           "查询失败",
		PAGE_ERROR:            "分页查询失败",

		JSON_ERROR: "json参数格式有误",
	}
)
