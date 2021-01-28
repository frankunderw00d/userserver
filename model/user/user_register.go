package user

type (
	// 注册请求
	RegisterRequest struct {
		PlatformID   int    `json:"platform_id"`   // 所属平台 ID
		Account      string `json:"account"`       // 账号
		Password     string `json:"password"`      // 秘密
		RegisterType int    `json:"register_type"` // 注册类型 0-游客 1-绑定用户
	}

	//// 注册响应
	//RegisterResponse struct {
	//	RegisterRequest
	//	Token string `json:"-"` // 账号绑定 token，内部分配唯一标识
	//}
)

const (
	AccountTableName = "dynamic_account"

	// 游客类型
	RegisterTypeTourist = 0
	// 绑定用户类型
	RegisterTypeCustomer = 1
)
