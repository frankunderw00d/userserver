package user

import (
	"baseservice/base/session"
	"encoding/json"
	"jarvis/base/network"
	loginModel "userserver/model/user"
)

type ()

const (
	// 上下文传递额外信息键
	ContextExtraSessionKey = "Session"
)

var ()

// 校验 Session 中间件函数
func (um *userModule) auth(ctx network.Context) {
	// 反序列化数据
	request := loginModel.AuthTypeRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		ctx.Done()
		return
	}

	// 校验 Session
	newSession, err := session.VerifySessionAndUpdate(request.Token, request.Session, request.SecretKey)
	if err != nil {
		printReplyError(ctx.ServerError(err))
		ctx.Done()
		return
	}

	// 存入上下文
	ctx.SetExtra(ContextExtraSessionKey, newSession)
}
