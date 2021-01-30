package user

import "baseservice/model/authenticate"

type (
	// 获取用户信息(除了用户 vip 等级,账号余额)请求
	GetUserInfoRequest struct {
		authenticate.Request
	}

	// 获取用户信息(除了用户 vip 等级,账号余额)响应
	GetUserInfoResponse struct {
		authenticate.Response
		AccountType       int    `json:"type"`                 // 账号类型 0-游客 1-绑定用户
		Platform          int    `json:"platform"`             // 所属平台
		Name              string `json:"name"`                 // 用户名
		Age               int    `json:"age"`                  // 用户年龄
		Sex               bool   `json:"sex"`                  // 用户性别
		HeadImage         int    `json:"head_image"`           // 用户头像序号
		Vip               int    `json:"vip"`                  // 用户 vip 等级
		GameBgMusicVolume int    `json:"game_bg_music_volume"` // 背景音乐音量
		GameEffectVolume  int    `json:"game_effect_volume"`   // 音效音量
		AccountBalance    int64  `json:"account_balance"`      // 账户余额(单位:分)
	}
)
