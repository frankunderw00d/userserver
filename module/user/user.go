package user

import (
	"baseservice/model/authenticate"
	"jarvis/base/network"
	"log"
)

type (
	// 模块定义
	userModule struct{}
)

const (
	// 模块名常量定义
	ModuleName = "User"
)

var (
	// 默认模块，由于需要同时充当观察者和模块，所以默认一个模块来支持不同的接口
	defaultUserModule *userModule
)

func init() {
	// 实例化默认模块
	defaultUserModule = &userModule{}
}

// 将默认模块声明为观察者
func NewObserver() network.Observer {
	return defaultUserModule
}

// 将默认模块声明为模块
func NewModule() network.Module {
	return defaultUserModule
}

// 模块要求实现函数: Name() string
func (um *userModule) Name() string {
	return ModuleName
}

// 模块要求实现函数: Route() map[string][]network.RouteHandleFunc
func (um *userModule) Route() map[string][]network.RouteHandleFunc {
	return map[string][]network.RouteHandleFunc{
		"register":             {um.register},                                        // 用户注册
		"login":                {um.login},                                           // 用户登录
		"getUserInfo":          {authenticate.Authenticate, um.getUserInfo},          // 获取用户信息(需要校验 Session )
		"updateUserInfo":       {authenticate.Authenticate, um.updateUserInfo},       // 更新用户信息(需要校验 Session )
		"updateAccountBalance": {authenticate.Authenticate, um.updateAccountBalance}, // 更新用户账户余额(需要校验 Session )
	}
}

// 观察者要求实现函数: ObserveConnect(string)
func (um *userModule) ObserveConnect(id string) {}

// 观察者要求实现函数: ObserveDisconnect(string)
func (um *userModule) ObserveDisconnect(id string) {}

// 观察者要求实现函数: InitiativeSend(network.Context)
func (um *userModule) InitiativeSend(ctx network.Context) {}

// 打印回复错误
func printReplyError(err error) {
	if err == nil {
		return
	}

	log.Printf("Reply error : %s", err.Error())
}
