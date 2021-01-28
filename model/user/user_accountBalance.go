package user

type (
	UpdateAccountBalanceRequest struct {
		AuthTypeRequest
		Amount   int64  `json:"amount"`   // 数额
		Describe string `json:"describe"` // 说明
	}

	UpdateAccountBalanceResponse struct {
		AuthTypeResponse
		AfterAmount int64 `json:"after_amount"` // 更新后数额
	}
)
