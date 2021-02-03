package user

import "baseservice/middleware/authenticate"

type (
	UpdateAccountBalanceRequest struct {
		authenticate.Request
		Amount     int64  `json:"amount"`      // 数额
		Describe   string `json:"describe"`    // 说明
		UpdateType int    `json:"update_type"` // 变动类型 0-充值 1-提现 2-游戏变动 3-活动奖励
	}

	UpdateAccountBalanceResponse struct {
		authenticate.Response
		AfterAmount int64 `json:"after_amount"` // 更新后数额
	}
)
