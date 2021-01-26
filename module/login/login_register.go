package login

import (
	"baseservice/model/platform"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jarvis/base/database"
	"jarvis/base/network"
	"jarvis/util/rand"
	"jarvis/util/regexp"
	loginModel "userserver/model/login"
)

// 注册
func (lm *loginModule) register(ctx network.Context) {
	request := loginModel.RegisterRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		return
	}

	response := &loginModel.RegisterResponse{}
	err := register(request, response)
	if err != nil {
		fmt.Printf("register error : %s", err.Error())
		printReplyError(ctx.ServerError(err))
		return
	}

	data, err := json.Marshal(&response)
	if err != nil {
		fmt.Printf("marshal response error : %s", err.Error())
		printReplyError(ctx.ServerError(err))
		return
	}

	printReplyError(ctx.Success(data))
}

func register(request loginModel.RegisterRequest, response *loginModel.RegisterResponse) error {
	// 获取 Redis 连接
	redisConn, err := database.GetRedisConn()
	if err != nil {
		return err
	}
	defer redisConn.Close()

	// 验证平台号
	if !platform.HExistsPlatformByID(fmt.Sprintf("%d", request.PlatformID)) {
		return errors.New("platform id doesn't exists")
	}

	// 绑定用户登录需要验证账号、秘密
	if request.RegisterType == loginModel.RegisterTypeCustomer {
		// 账号要求 6-18位，只允许字母数字，不允许数字开头
		if !regexp.Match("^[a-zA-Z]+[a-zA-Z0-9]{5,17}$", request.Account) {
			return errors.New("require account length must between 6 - 18")
		}

		// 密码要求 6-18位，只允许字母数字
		if !regexp.Match("^[a-zA-Z0-9]{6,18}$", request.Password) {
			return errors.New("require password length must between 6 - 18")
		}
	} else {
		// 随机分配账号密码
		request.Account = rand.RandomString(10, rand.SeedUCL, rand.SeedLCL) + rand.RandomString(8, rand.SeedNum)
		request.Password = rand.RandomString(18, rand.SeedUCL, rand.SeedLCL, rand.SeedNum)
	}

	// 生成唯一 token
	token := rand.RandomString(16)

	// 获取 MySQL 连接
	mysqlConn, err := database.GetMySQLConn()
	if err != nil {
		return err
	}
	defer mysqlConn.Close()

	if _, err := mysqlConn.ExecContext(context.Background(),
		"insert into `dynamic_account`(token,account,password,`type`,platform)values(?,?,?,?,?);",
		token, request.Account, request.Password, request.RegisterType, request.PlatformID,
	); err != nil {
		return err
	}

	response.RegisterRequest = request
	response.Token = token

	return nil
}
