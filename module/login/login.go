package login

import (
	"context"
	"jarvis/base/database"
	"jarvis/base/network"
	"log"
	"sync"
	globalModel "userserver/model/global"
)

type (
	// 模块定义
	loginModule struct {
		userMap     map[string]bool // 用户映射表，用于用户是否校验通过
		userMapMute sync.Mutex      // 用户服务必须对用户的数据进行锁定处理，防止并发问题
	}
)

const (
	// 模块名常量定义
	ModuleName = "Login"
)

var (
	// 默认模块，由于需要同时充当观察者和模块，所以默认一个模块来支持不同的接口
	defaultLoginModule *loginModule
	// 全局平台列表
	globalPlatformList globalModel.PlatformList
)

func init() {
	// 实例化默认模块
	defaultLoginModule = &loginModule{
		userMap:     make(map[string]bool),
		userMapMute: sync.Mutex{},
	}

	// 实例化全局平台列表
	globalPlatformList = make(globalModel.PlatformList, 0)
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
		"authenticate": {lm.authenticate},   // 校验函数
		"login":        {lm.auth, lm.login}, // 登录函数，前置校验检查中间件
	}
}

// 观察者要求实现函数: ObserveConnect(string)
func (lm *loginModule) ObserveConnect(id string) {
	lm.userMapMute.Lock()
	defer lm.userMapMute.Unlock()

	if _, exist := lm.userMap[id]; exist {
		log.Printf("[%s] module observe that [%s] exist", ModuleName, id)
		return
	}

	// false - 等待校验 ， true - 校验通过
	lm.userMap[id] = false
}

// 观察者要求实现函数: ObserveDisconnect(string)
func (lm *loginModule) ObserveDisconnect(id string) {
	lm.userMapMute.Lock()
	defer lm.userMapMute.Unlock()

	if _, exist := lm.userMap[id]; !exist {
		log.Printf("[%s] module observe that [%s] unexist", ModuleName, id)
		return
	}

	delete(lm.userMap, id)
}

// 打印回复错误
func printReplyError(err error) {
	if err == nil {
		return
	}

	log.Printf("Reply error : %s", err.Error())
}

// 加载平台信息
func LoadPlatformInfo() {
	log.Println("======================================================================")
	log.Printf("Start loading platform information from `jarvis`.`static_platform`")
	conn, err := database.GetMySQLConn()
	if err != nil {
		log.Fatalf("database.GetMySQLConn error : %s", err.Error())
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("database.GetMySQLConn close error : %s", err.Error())
			return
		}
	}()

	platformList := make(globalModel.PlatformList, 0)

	rows, err := conn.QueryContext(context.Background(), platformList.QueryOrder())
	if err != nil {
		log.Fatalf("database.GetMySQLConn query error : %s", err.Error())
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatalf("database.GetMySQLConn rows close error : %s", err.Error())
			return
		}
	}()

	for rows.Next() {
		platform := globalModel.Platform{}
		err := rows.Scan(&platform.ID, &platform.Name, &platform.Link, &platform.Owner, &platform.CreateAt, &platform.UpdateAt)
		if err != nil {
			log.Printf("database.GetMySQLConn query rows scan error : %s", err.Error())
			return
		}
		platformList = append(platformList, platform)
	}

	log.Printf("Now we having %d platform:", len(platformList))
	if len(platformList) > 0 {
		globalPlatformList = platformList
	}
	for _, platform := range globalPlatformList {
		log.Printf("Platform : %+v", platform)
	}
	log.Printf("Finish loading platform information from `jarvis`.`static_platform`")
	log.Println("======================================================================")
}
