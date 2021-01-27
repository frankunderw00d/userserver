package login

import (
	"baseservice/base/basic"
	"baseservice/model/user"
	"context"
	"encoding/json"
	"fmt"
	"jarvis/base/database"
	"jarvis/base/network"
	uRand "jarvis/util/rand"
	"log"
	loginModel "userserver/model/login"
)

var ()

const (
	UsersInfoKey                      = "UsersInfo"
	UserInfoField basic.ComposeString = "User:"
)

var ()

// 登录
func (lm *loginModule) login(ctx network.Context) {
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

func login(request loginModel.LoginRequest, response *loginModel.LoginResponse) error {
	// 获取 MySQL 连接
	mysqlConn, err := database.GetMySQLConn()
	if err != nil {
		return err
	}
	defer mysqlConn.Close()

	// 获取用户信息
	freshUser := user.FreshUser()
	row := mysqlConn.QueryRowContext(context.Background(), "select id,token,account,type,platform from `jarvis`.`dynamic_account` where account = ? and password = ?",
		request.Account, request.Password)
	err = row.Scan(&freshUser.Account.ID, &freshUser.Account.Token, &freshUser.Account.Account, &freshUser.Account.AccountType, &freshUser.Account.Platform)
	if err != nil {
		return err
	}

	log.Printf("%+v", freshUser)

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

	response.Token = freshUser.Account.Token
	response.Session = session

	return nil
}

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
