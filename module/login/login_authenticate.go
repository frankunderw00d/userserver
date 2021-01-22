package login

import (
	"encoding/json"
	"jarvis/base/network"
	loginModel "userserver/model/login"
)

var (
	fakeUserAuthInfo = map[string]string{
		"frank":  "12345678",
		"frank1": "123456781",
		"frank2": "123456782",
		"frank3": "123456783",
		"frank4": "123456784",
		"frank5": "123456785",
		"frank6": "123456786",
	}
)

// 校验接口
func (lm *loginModule) authenticate(ctx network.Context) {
	// 加锁，防止竞态
	lm.userMapMute.Lock()
	defer lm.userMapMute.Unlock()

	// 反序列化
	request := loginModel.AuthenticateReq{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.Error(err))
		return
	}

	// 校验信息
	if request.Account == "" {
		printReplyError(ctx.Error(ctx.Reply(ctx.Request().ID, ctx.Request().Reply, []byte("account can't be nil"))))
		return
	}

	token, exist := fakeUserAuthInfo[request.Account]
	if !exist {
		printReplyError(ctx.Error(ctx.Reply(ctx.Request().ID, ctx.Request().Reply, []byte("account doesn't exist"))))
		return
	}

	if request.Token != token {
		printReplyError(ctx.Error(ctx.Reply(ctx.Request().ID, ctx.Request().Reply, []byte("token error"))))
		return
	}

	// 确认连接的真实性，如果中途意外下线，即 lm.userMap 中没有，此消息将会不处理
	if _, exist := lm.userMap[ctx.Request().ID]; !exist {
		return
	}

	lm.userMap[ctx.Request().ID] = true
	printReplyError(ctx.Reply(ctx.Request().ID, ctx.Request().Reply, []byte("authenticate success")))
}

// 请求是否通过校验中间件
func (lm *loginModule) auth(ctx network.Context) {
	// 加锁，防止竞态
	lm.userMapMute.Lock()
	defer lm.userMapMute.Unlock()

	// 确认请求的有效性，ID 不存在 或 未通过校验，都会引发中断调用链
	if pass, exist := lm.userMap[ctx.Request().ID]; !exist || !pass {
		printReplyError(ctx.Error(ctx.Reply(ctx.Request().ID, ctx.Request().Reply, []byte("authenticate failed,please verify first!"))))
		// 中断调用链
		ctx.Done()
		return
	}
}
