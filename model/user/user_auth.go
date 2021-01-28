package user

type (
	// 前置验证类型的请求
	AuthTypeRequest struct {
		Token     string `json:"token"`      // 账号唯一标识
		Session   string `json:"session"`    // 会话标识
		SecretKey string `json:"secret_key"` // 加密 key
	}

	// 前置验证类型的响应
	AuthTypeResponse struct {
		Session string `json:"session"` // 更新会话标识
	}
)
