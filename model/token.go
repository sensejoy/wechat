package model

import (
	"wechat/model/dao"
	"wechat/util"
)

const platformKey = "wechat:platform:app_token"
const miniProgramKey = "wechat:miniprogram:app_token"

func SetAccessToken(token string, kind int) {
	if len(token) == 0 {
		util.Logger.Error("invalid token")
	}
	conn := dao.RedisPool.Get()
	var key string
	if kind == util.OfficialAccount {
		key = platformKey
	} else if kind == util.MiniProgram {
		key = miniProgramKey
	} else {
		util.Logger.Error("invalid token type")
		return
	}

	defer conn.Close()
	_, err := conn.Do("set", key, token)
	if err == nil {
		util.Logger.Info("update wechat access token successfully")
	} else {
		util.Logger.Error("failed to update wechat access token")
	}
}

func GetAccessToken(kind int) string {
	conn := dao.RedisPool.Get()
	defer conn.Close()
	var key string

	if kind == util.OfficialAccount {
		key = platformKey
	} else if kind == util.MiniProgram {
		key = miniProgramKey
	} else {
		util.Logger.Error("invalid token type")
		return ""
	}
	token, err := conn.Do("get", key)
	if err == nil {
		return string(token.([]byte))
	} else {
		return ""
	}
}
