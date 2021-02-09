package user

import (
	"baseservice/model/platform"
	"baseservice/model/user"
	"encoding/json"
	"errors"
	"fmt"
	"jarvis/base/log"
	"jarvis/base/network"
	"jarvis/util/rand"
	"jarvis/util/regexp"
	loginModel "userserver/model/user"
)

// 注册
func (um *userModule) register(ctx network.Context) {
	// 反序列化数据
	request := loginModel.RegisterRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		return
	}

	// 调用业务函数
	err := register(request)
	if err != nil {
		log.ErrorF("register error : %s", err.Error())
		printReplyError(ctx.ServerError(err))
		return
	}

	// 返回响应
	printReplyError(ctx.Success([]byte("Register succeed")))
}

func register(request loginModel.RegisterRequest) error {
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

	freshUser := user.FreshUser()
	freshUser.Account.Token = token
	freshUser.Account.Account = request.Account
	freshUser.Account.Password = request.Password
	freshUser.Account.AccountType = request.RegisterType
	freshUser.Account.Platform = request.PlatformID
	freshUser.Info.Name = rand.RandomString(10)
	if err := freshUser.Store(); err != nil {
		return err
	}

	return nil
}
