package model

import (
	"wechat/model/dao"
	"wechat/util"
)

const redisKey = "wechat:app_token"

func SetAccessToken(token string) {
	if len(token) == 0 {
		util.Logger.Error("invalid token")
	}
	conn := dao.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("set", redisKey, token)
	if err == nil {
		util.Logger.Info("update wechat access token successfully")
	} else {
		util.Logger.Error("failed to update wechat access token")
	}
}

func GetAccessToken() string {
	conn := dao.RedisPool.Get()
	defer conn.Close()
	token, err := conn.Do("get", redisKey)
	if err == nil {
		return string(token.([]byte))
	} else {
		return ""
	}
}
