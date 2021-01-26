package login

import (
	"jarvis/base/network"
)

var ()

func (lm *loginModule) login(ctx network.Context) {
	//// 加锁，防止竞态
	//lm.userMapMute.Lock()
	//defer lm.userMapMute.Unlock()
	//
	//// 反序列化
	//request := loginModel.LoginReq{}
	//if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
	//	printReplyError(ctx.Error(err))
	//	return
	//}
	//
	//// 校验信息
	//if request.Account == "" {
	//	printReplyError(ctx.Error(ctx.Reply(ctx.Request().ID, ctx.Request().Reply, []byte("account can't be nil"))))
	//	return
	//}
	//
	//password, exist := fakeUserLoginInfo[request.Account]
	//if !exist {
	//	printReplyError(ctx.Error(ctx.Reply(ctx.Request().ID, ctx.Request().Reply, []byte("account doesn't exist"))))
	//	return
	//}
	//
	//if request.Password != password {
	//	printReplyError(ctx.Error(ctx.Reply(ctx.Request().ID, ctx.Request().Reply, []byte("token error"))))
	//	return
	//}
	//
	//// 确认连接的真实性，如果中途意外下线，即 lm.userMap 中没有，此消息将会不处理
	//if _, exist := lm.userMap[ctx.Request().ID]; !exist {
	//	return
	//}
	//
	//// 发送响应
	//response := loginModel.LoginRsp{
	//	Session: rand.RandomString(32),
	//}
	//rspData, err := json.Marshal(&response)
	//if err != nil {
	//	printReplyError(ctx.Error(err))
	//	return
	//}
	//
	//printReplyError(ctx.Reply(ctx.Request().ID, ctx.Request().Reply, rspData))
}
