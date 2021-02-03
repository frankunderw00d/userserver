package user

import "baseservice/middleware/authenticate"

type (
	// 更新用户信息(除了用户 vip 等级,账号余额)请求
	UpdateRequest struct {
		authenticate.Request
		Name              string `json:"name"`                 // 用户名
		Age               int    `json:"age"`                  // 用户年龄
		Sex               bool   `json:"sex"`                  // 用户性别
		HeadImage         int    `json:"head_image"`           // 用户头像序号
		GameBgMusicVolume int    `json:"game_bg_music_volume"` // 背景音乐音量
		GameEffectVolume  int    `json:"game_effect_volume"`   // 音效音量
	}

	// 更新用户信息(除了用户 vip 等级,账号余额)响应
	UpdateResponse struct {
		GetUserInfoResponse
	}
)
