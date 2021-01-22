package login

type (
	AuthenticateReq struct {
		Account string `json:"account"` // 账号
		Token   string `json:"token"`   // 密钥
	}

	LoginReq struct {
		Account  string `json:"account"`  // 账号
		Password string `json:"password"` // 密码
	}

	LoginRsp struct {
		Session string `json:"session"` // 临时会话密钥
	}
)
