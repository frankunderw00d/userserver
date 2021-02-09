package user

import (
	"baseservice/middleware/authenticate"
	"baseservice/model/user"
	"encoding/json"
	"jarvis/base/database/redis"
	"jarvis/base/log"
	"jarvis/base/network"
	"time"
	userModel "userserver/model/user"
)

// 用户下线，清理缓存信息，记录用户下线时间
func (um *userModule) logout(ctx network.Context) {
	// 反序列化数据
	request := userModel.LogoutRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		return
	}

	// 实例化响应
	response := &userModel.LogoutResponse{}
	// 调用函数
	err := logout(request, response)
	if err != nil {
		log.ErrorF("get user info error : %s", err.Error())
		printReplyError(ctx.ServerError(err))
		return
	}

	newSession := ctx.Extra(authenticate.ContextExtraSessionKey, "")
	response.Session = newSession.(string)

	// 序列化响应
	data, err := json.Marshal(response)
	if err != nil {
		log.ErrorF("marshal response error : %s", err.Error())
		printReplyError(ctx.ServerError(err))
		return
	}

	// 返回响应
	printReplyError(ctx.Success(data))
}

func logout(request userModel.LogoutRequest, response *userModel.LogoutResponse) error {
	// 1.用户的 Session 不用删除，下次用户可以用 token + session 的方式自动登录
	// 2.删除用户在 redis 中的信息
	if _, err := redis.HDel(UsersInfoKey, UserInfoField.Compose(request.Token)); err != nil {
		return err
	}

	// 3.记录用户登出时间
	// 加载用户信息
	freshUser := user.FreshUser()
	if err := freshUser.LoadInfoByToken(request.Token); err != nil {
		return err
	}

	// 存入登出时间
	freshUser.Account.LastLogout = time.Now()
	if err := freshUser.Account.StoreLogoutTime(); err != nil {
		return err
	}

	return nil
}
