package user

import (
	"baseservice/base/basic"
	"baseservice/model/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jarvis/base/database"
	"jarvis/base/network"
	uRand "jarvis/util/rand"
	loginModel "userserver/model/user"
)

var ()

const (
	UsersInfoKey                      = "UsersInfo"
	UserInfoField basic.ComposeString = "User:"
)

var ()

// 登录
// 登录函数为所有函数的前提
// 因此登录函数强制刷新 redis 账号绑定的 Session 和 用户信息
func (um *userModule) login(ctx network.Context) {
	// 反序列化数据
	request := loginModel.LoginRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		return
	}

	// 实例化响应
	response := &loginModel.LoginResponse{}
	// 调用函数
	err := login(request, response)
	if err != nil {
		fmt.Printf("user error : %s", err.Error())
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

func login(request loginModel.LoginRequest, response *loginModel.LoginResponse) error {
	// 获取 MySQL 连接
	mysqlConn, err := database.GetMySQLConn()
	if err != nil {
		return err
	}
	defer mysqlConn.Close()

	// 验证账号合法性
	var exist int
	row := mysqlConn.QueryRowContext(context.Background(), "select count(id) from `jarvis`.`dynamic_account` where account = ? and password = ?",
		request.Account, request.Password)
	err = row.Scan(&exist)
	if err != nil {
		return err
	}

	// 不存在
	if exist < 1 {
		return errors.New("account or password wrong")
	}

	// 加载用户信息
	freshUser := user.FreshUser()
	if err := freshUser.LoadInfoByAccountAndPassword(request.Account, request.Password); err != nil {
		return err
	}

	// 生成随机 Session
	session := uRand.RandomString(8)

	// 存入用户信息到 redis
	if err := SetUserInfoToRedis(freshUser); err != nil {
		return err
	}

	// 存入用户 Session 到 redis
	if err := SetSession(freshUser.Account.Token, session); err != nil {
		return err
	}

	// 返回
	response.Token = freshUser.Account.Token
	response.Session = session

	return nil
}

// 将用户信息存入 redis
func SetUserInfoToRedis(u user.User) error {
	// 获取 Redis 连接
	redisConn, err := database.GetRedisConn()
	if err != nil {
		return err
	}
	defer redisConn.Close()

	userData, err := json.Marshal(&u)
	if err != nil {
		return err
	}

	_, err = redisConn.Do("hset", UsersInfoKey, UserInfoField.Compose(u.Account.Token), string(userData))
	return err
}
