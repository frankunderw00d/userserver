package user

import "baseservice/model/authenticate"

type (
	UpdateAccountBalanceRequest struct {
		authenticate.Request
		Amount   int64  `json:"amount"`   // 数额
		Describe string `json:"describe"` // 说明
	}

	UpdateAccountBalanceResponse struct {
		authenticate.Response
		AfterAmount int64 `json:"after_amount"` // 更新后数额
	}
)
