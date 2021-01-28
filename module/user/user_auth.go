package user

import (
	"baseservice/base/basic"
	"encoding/json"
	redisGo "github.com/gomodule/redigo/redis"
	"jarvis/base/database"
	"jarvis/base/network"
	uRand "jarvis/util/rand"
	"strings"
	"time"
	loginModel "userserver/model/user"
)

type ()

const (
	// 用户 Session 键
	UsersSessionKey = "UsersSession"
	// 用户 Session 列键
	UsersSessionField basic.ComposeString = "User:"
	// 上下文传递额外信息键
	ContextExtraSessionKey = "Session"
)

var ()

func (um *userModule) auth(ctx network.Context) {
	// 反序列化数据
	request := loginModel.AuthTypeRequest{}
	if err := json.Unmarshal(ctx.Request().Data, &request); err != nil {
		printReplyError(ctx.ServerError(err))
		ctx.Done()
		return
	}

	// 根据 request.Token 取得 Session
	redisSession, err := GetSession(request.Token)
	if err != nil {
		printReplyError(ctx.ServerError(err))
		ctx.Done()
		return
	}

	// 核对 Session 和 secretKey
	if redisSession != request.Session {
		printReplyError(ctx.BadRequest("session wrong"))
		ctx.Done()
		return
	}
	if basic.EncryptSecretKey(request.Token, redisSession) != request.SecretKey {
		printReplyError(ctx.BadRequest("secretKey wrong"))
		ctx.Done()
		return
	}

	// 核对 Session 是否超时，超时则更换，存入 Redis 且返回给用户
	timeout, err := CheckSessionTimeout(request.Token)
	if err != nil {
		printReplyError(ctx.ServerError(err))
		ctx.Done()
		return
	}

	if timeout {
		// 生成随机 Session
		session := uRand.RandomString(8)
		err := SetSession(request.Token, session)
		if err != nil {
			printReplyError(ctx.ServerError(err))
			ctx.Done()
			return
		}

		// 在调用链中传递额外信息
		ctx.SetExtra(ContextExtraSessionKey, session)
	}
}

// 设置 Session
func SetSession(token, session string) error {
	// 获取 Redis 连接
	redisConn, err := database.GetRedisConn()
	if err != nil {
		return err
	}
	defer redisConn.Close()

	now := time.Now().Format("20060102150405")
	_, err = redisConn.Do("hset", UsersSessionKey, UsersSessionField.Compose(token), session+":"+now)
	if err != nil {
		return err
	}

	return nil
}

// 获取 Session
func GetSession(token string) (string, error) {
	// 获取 Redis 连接
	redisConn, err := database.GetRedisConn()
	if err != nil {
		return "", err
	}
	defer redisConn.Close()

	v, err := redisGo.String(redisConn.Do("hget", UsersSessionKey, UsersSessionField.Compose(token)))
	if err != nil {
		return "", err
	}

	return strings.SplitN(v, ":", 2)[0], nil
}

// Session 是否超时，默认15分钟
func CheckSessionTimeout(token string) (bool, error) {
	// 获取 Redis 连接
	redisConn, err := database.GetRedisConn()
	if err != nil {
		return false, err
	}
	defer redisConn.Close()

	v, err := redisGo.String(redisConn.Do("hget", UsersSessionKey, UsersSessionField.Compose(token)))
	if err != nil {
		return false, err
	}

	p, err := time.ParseInLocation("20060102150405", strings.SplitN(v, ":", 2)[1], time.Local)
	if err != nil {
		return false, err
	}

	if time.Now().Sub(p).Minutes() >= 15 {
		return true, nil
	}

	return false, nil
}
