/*
	HashMap UsersSession User:Token Session:Time
*/
package login

import (
	"baseservice/base/basic"
	redisGo "github.com/gomodule/redigo/redis"
	"jarvis/base/database"
	"jarvis/base/network"
	"strings"
	"time"
)

type ()

const (
	UsersSessionKey                       = "UsersSession"
	UsersSessionField basic.ComposeString = "User:"
)

var ()

func (lm *loginModule) auth(ctx network.Context) {

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

	if time.Now().Sub(p).Minutes() >= -15 {
		return true, nil
	}

	return false, nil
}