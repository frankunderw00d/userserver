package login

import (
	"jarvis/base/network"
	"log"
)

type (
	// 模块定义
	loginModule struct{}
)

const (
	// 模块名常量定义
	ModuleName = "Login"
)

var (
	// 默认模块，由于需要同时充当观察者和模块，所以默认一个模块来支持不同的接口
	defaultLoginModule *loginModule
)

func init() {
	// 实例化默认模块
	defaultLoginModule = &loginModule{}
}

// 将默认模块声明为观察者
func NewObserver() network.Observer {
	return defaultLoginModule
}

// 将默认模块声明为模块
func NewModule() network.Module {
	return defaultLoginModule
}

// 模块要求实现函数: Name() string
func (lm *loginModule) Name() string {
	return ModuleName
}

// 模块要求实现函数: Route() map[string][]network.RouteHandleFunc
func (lm *loginModule) Route() map[string][]network.RouteHandleFunc {
	return map[string][]network.RouteHandleFunc{
		"register": {lm.register},       // 注册
		"login":    {lm.auth, lm.login}, // 登录函数，前置校验检查中间件
	}
}

// 观察者要求实现函数: ObserveConnect(string)
func (lm *loginModule) ObserveConnect(id string) {}

// 观察者要求实现函数: ObserveDisconnect(string)
func (lm *loginModule) ObserveDisconnect(id string) {}

// 打印回复错误
func printReplyError(err error) {
	if err == nil {
		return
	}

	log.Printf("Reply error : %s", err.Error())
}
