package user

import (
	"baseservice/base/basic"
	"baseservice/model/user"
	"encoding/json"
	"errors"
	"fmt"
	"jarvis/base/database"
	"jarvis/base/network"
	loginModel "userserver/model/user"
)

const (
	// 用户信息分布式锁名
	InfoUpdateDistributedLock basic.ComposeString = "UserInfoUpdateDisLock:"
	// 更新用户信息分布式加锁失败文字
	ErrInfoUpdateDistributedLockText = "information update distributed lock failure"
)

var (
	// 更新用户信息分布式加锁失败
	ErrInfoUpdateDistributedLock = errors.New(ErrInfoUpdateDistributedLockText)
)

// 更新用户信息(除了用户 vip 等级,账号余额)
func (um *userModule) updateUserInfo(ctx network.Context) {
	// 反序列化数据
	request := loginModel.UpdateRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		return
	}

	// 实例化响应
	response := &loginModel.UpdateResponse{}
	// 调用函数
	err := updateUserInfo(request, response)
	if err != nil {
		fmt.Printf("update user info error : %s", err.Error())
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

func updateUserInfo(request loginModel.UpdateRequest, response *loginModel.UpdateResponse) error {
	// 获取 redis 链接
	redisConn, err := database.GetRedisConn()
	if err != nil {
		return err
	}
	defer redisConn.Close()

	// 新建 redis 分布式锁
	redisLock := database.NewRedisLock()
	if err := redisLock.Initialize(); err != nil {
		return err
	}
	defer redisLock.Close()

	// 上锁
	if !redisLock.UntilLock(InfoUpdateDistributedLock.Compose(request.Token),10) {
		return ErrInfoUpdateDistributedLock
	}
	defer redisLock.Unlock(InfoUpdateDistributedLock.Compose(request.Token))

	// 加载用户信息
	freshUser := user.FreshUser()
	if err := freshUser.LoadInfoByToken(request.Token); err != nil {
		return err
	}

	// 赋予新值
	freshUser.Info.Name = request.Name
	freshUser.Info.Age = request.Age
	freshUser.Info.Sex = request.Sex
	freshUser.Info.HeadImage = request.HeadImage
	freshUser.Info.GameBgMusicVolume = request.GameBgMusicVolume
	freshUser.Info.GameEffectVolume = request.GameEffectVolume

	// 更新到 MySQL 中
	if err := freshUser.Info.Update(); err != nil {
		return err
	}

	// 存入用户信息到 redis
	if err := SetUserInfoToRedis(freshUser); err != nil {
		return err
	}

	// 返回
	response.AccountType = freshUser.Account.AccountType
	response.Platform = freshUser.Account.Platform
	response.Name = freshUser.Info.Name
	response.Age = freshUser.Info.Age
	response.Sex = freshUser.Info.Sex
	response.HeadImage = freshUser.Info.HeadImage
	response.Vip = freshUser.Info.Vip
	response.GameBgMusicVolume = freshUser.Info.GameBgMusicVolume
	response.GameEffectVolume = freshUser.Info.GameEffectVolume
	response.AccountBalance = freshUser.Info.AccountBalance

	return nil
}
