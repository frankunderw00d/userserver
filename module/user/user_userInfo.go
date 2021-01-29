package user

import (
	"baseservice/model/user"
	"encoding/json"
	"fmt"
	redisGo "github.com/gomodule/redigo/redis"
	"jarvis/base/database"
	"jarvis/base/network"
	loginModel "userserver/model/user"
)

// 获取用户信息
func (um *userModule) getUserInfo(ctx network.Context) {
	// 反序列化数据
	request := loginModel.GetUserInfoRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		return
	}

	// 实例化响应
	response := &loginModel.GetUserInfoResponse{}
	// 调用函数
	err := getUserInfo(request, response)
	if err != nil {
		fmt.Printf("get user info error : %s", err.Error())
		printReplyError(ctx.ServerError(err))
		return
	}

	newSession := ctx.Extra(ContextExtraSessionKey, "")
	response.Session = newSession.(string)

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

func getUserInfo(request loginModel.GetUserInfoRequest, response *loginModel.GetUserInfoResponse) error {
	// 从 redis 读取用户信息，不需要加锁
	// 加载用户信息
	u, err := GetUserInfoFromRedis(request.Token)
	if err != nil {
		return err
	}

	// 返回
	response.AccountType = u.Account.AccountType
	response.Platform = u.Account.Platform
	response.Name = u.Info.Name
	response.Age = u.Info.Age
	response.Sex = u.Info.Sex
	response.HeadImage = u.Info.HeadImage
	response.Vip = u.Info.Vip
	response.GameBgMusicVolume = u.Info.GameBgMusicVolume
	response.GameEffectVolume = u.Info.GameEffectVolume
	response.AccountBalance = u.Info.AccountBalance

	return nil
}

// 根据 token 从 redis 中获取用户数据
func GetUserInfoFromRedis(token string) (user.User, error) {
	// 获取 Redis 连接
	redisConn, err := database.GetRedisConn()
	if err != nil {
		return user.User{}, err
	}
	defer redisConn.Close()

	infoStr, err := redisGo.String(redisConn.Do("hget", UsersInfoKey, UserInfoField.Compose(token)))
	if err != nil {
		return user.User{}, err
	}

	u := user.User{}

	return u, json.Unmarshal([]byte(infoStr), &u)
}
