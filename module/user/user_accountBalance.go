package user

import (
	"baseservice/base/basic"
	"baseservice/model/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"jarvis/base/database"
	"jarvis/base/network"
	"time"
	loginModel "userserver/model/user"
)

const (
	UserAccountBalanceUpdateDistributedLock basic.ComposeString = "UserAccountBalanceUpdateDisLock:"
)

func (um *userModule) updateAccountBalance(ctx network.Context) {
	// 反序列化数据
	request := loginModel.UpdateAccountBalanceRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		return
	}

	// 实例化响应
	response := &loginModel.UpdateAccountBalanceResponse{}
	// 调用函数
	err := updateAccountBalance(request, response)
	if err != nil {
		fmt.Printf("user error : %s", err.Error())
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

func updateAccountBalance(request loginModel.UpdateAccountBalanceRequest, response *loginModel.UpdateAccountBalanceResponse) error {
	// 获取 redis 链接
	redisConn, err := database.GetRedisConn()
	if err != nil {
		return err
	}
	defer redisConn.Close()

	// 新建 redis 分布式锁
	redisLock := database.NewRedisLock(redisConn, 10)
	defer redisLock.Close()

	// 上锁
	if !redisLock.Lock(UserAccountBalanceUpdateDistributedLock.Compose(request.Token), 7) {
		return errors.New("distributed lock failure")
	}
	defer redisLock.UnLock(UserAccountBalanceUpdateDistributedLock.Compose(request.Token))

	// 加载用户信息
	freshUser := user.FreshUser()
	if err := freshUser.LoadInfoByToken(request.Token); err != nil {
		return err
	}

	// 赋予新值
	freshUser.Info.AccountBalance = freshUser.Info.AccountBalance + request.Amount

	// 更新到 MySQL 中
	if err := freshUser.Info.UpdateAccountBalance(); err != nil {
		return err
	}

	// 存入用户信息到 redis
	if err := SetUserInfoToRedis(freshUser); err != nil {
		return err
	}

	// 将账户更改记录存入 mongo
	collection, err := database.GetMongoConn("dynamic_user_account_balance_update_record")
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.Background(), bson.M{"amount": request.Amount, "time": time.Now().Format("2006-01-02 15:04:05"), "describe": request.Describe})
	if err != nil {
		return err
	}

	// 返回
	response.AfterAmount = freshUser.Info.AccountBalance

	return nil
}
