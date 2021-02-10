package main

import (
	"baseservice/middleware/traceRecord"
	"io"
	"jarvis/base/database"
	"jarvis/base/database/redis"
	"jarvis/base/log"
	"jarvis/base/network"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userserver/module/user"
	"userserver/utils/logHook"
)

const (
	// 自定义服务最大连接数量
	CustomMaxConnection = 5000
	// 自定义服务消息管道大小
	CustomIntoStreamSize = 10000
	// Socket 监听地址
	SocketListenAddress = ":8080"
	// WebSocket 监听地址
	WebSocketListenAddress = ":8081"
	// gRPC 监听地址
	GRPCListenAddress = ":8082"
	// 远程日志聚合地址
	LogRemoteAddress = ":10000"
)

var (
	service network.Service
	lh      io.WriteCloser
)

func init() {
	// 新建远程日志钩子
	nlh, err := logHook.NewRemoteHook(LogRemoteAddress)
	if err != nil {
		log.FatalF("New remote hook error : %s", err.Error())
		return
	}
	lh = nlh
	log.SetHook(lh)

	// 实例化服务
	service = network.NewService(
		CustomMaxConnection,
		CustomIntoStreamSize,
	)

	// 初始化 MySQL
	if err := database.InitializeMySQL(
		"frank",
		"frank123",
		"mysql-service",
		3306,
		"jarvis",
	); err != nil {
		log.FatalF("Initialize MySQL error : %s", err.Error())
		return
	}

	// 设置 MySQL
	database.SetUpMySQL(time.Minute*time.Duration(5), 10, 5000)

	// 初始化 Redis
	redis.InitializeRedis(time.Minute*time.Duration(5), 10, 5000, "redis-service", 6379, "frank123")

	// 初始化 Mongo
	if err := database.InitializeMongo(
		"frank",
		"frank123",
		"jarvis",
		"mongo-service",
		27017, time.Minute*time.Duration(5), 5000); err != nil {
		log.FatalF("Initialize Mongo error : %s", err.Error())
		return
	}
}

func main() {
	// 1.添加全局中间件
	if err := service.UseMiddleware(logMiddleware, traceRecord.TraceRecord); err != nil {
		log.ErrorF("Use middleware error : %s", err)
		return
	}

	// 2.注册观察者
	if err := service.RegisterObserver(user.NewObserver()); err != nil {
		log.ErrorF("Register observer error : %s", err)
		return
	}

	// 3.注册模块
	if err := service.RegisterModule(user.NewModule()); err != nil {
		log.ErrorF("Register module error : %s", err)
		return
	}

	// 4.启动
	if err := service.Run(
		network.NewSocketGate(SocketListenAddress),       // Socket 入口
		network.NewWebSocketGate(WebSocketListenAddress), // WebSocket 入口
		network.NewGRPCGate(GRPCListenAddress),           // gRPC 入口
	); err != nil {
		log.ErrorF("Register observer error : %s", err)
		return
	}

	// 5.监听系统信号
	monitorSystemSignal()
}

// 监听系统信号
// kill -SIGQUIT [进程号] : 杀死当前进程
func monitorSystemSignal() {
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGQUIT)
	select {
	case <-sc:
		_ = lh.Close()
		log.InfoF("Done")
	}
}

// 打印中间件
func logMiddleware(ctx network.Context) {
	log.InfoF("%+v", ctx.Request())
}
