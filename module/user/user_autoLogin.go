package user

import (
	bSession "baseservice/common/session"
	"baseservice/model/user"
	"encoding/json"
	"fmt"
	"jarvis/base/network"
	uRand "jarvis/util/rand"
	"time"
	userModel "userserver/model/user"
)

// 自动登录
// 自动登录和登录函数为所有函数的前提
// 因此登录函数强制刷新 redis 账号绑定的 Session 和 用户信息
func (um *userModule) autoLogin(ctx network.Context) {
	// 反序列化数据
	request := userModel.AutoLoginRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		return
	}

	// 实例化响应
	response := &userModel.AutoLoginResponse{}
	// 调用函数
	err := autoLogin(request, response)
	if err != nil {
		fmt.Printf("login error : %s", err.Error())
		printReplyError(ctx.ServerError(err))
		return
	}

	// 序列化响应
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("marshal response error : %s", err.Error())
		printReplyError(ctx.ServerError(err))
		return
	}

	// 返回响应
	printReplyError(ctx.Success(data))
}

func autoLogin(request userModel.AutoLoginRequest, response *userModel.AutoLoginResponse) error {
	// 加载用户信息
	freshUser := user.FreshUser()
	if err := freshUser.LoadInfoByToken(request.Token); err != nil {
		return err
	}

	// 存入登录时间
	freshUser.Account.LastLogin = time.Now()
	if err := freshUser.Account.StoreLoginTime(); err != nil {
		return err
	}

	// 生成随机 Session
	session := uRand.RandomString(8)

	// 存入用户信息到 redis
	if err := SetUserInfoToRedis(freshUser); err != nil {
		return err
	}

	// 存入用户 Session 到 redis
	if err := bSession.SetSession(freshUser.Account.Token, session); err != nil {
		return err
	}

	// 返回
	response.Session = session

	return nil
}
