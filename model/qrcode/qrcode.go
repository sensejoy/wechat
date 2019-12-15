package qrcode

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"wechat/model"
	"wechat/model/dao"
	"wechat/util"
)

type qrCodePostInfo struct {
	Expire int        `json:"expire_seconds"`
	Name   string     `json:"action_name"`
	Info   actionInfo `json:"action_info"`
}

type actionInfo struct {
	Scene sceneInfo `json:"scene"`
}

type sceneInfo struct {
	Id  int    `json:"scene_id"`
	Str string `json:"scene_str"`
}

/**
 * desc 创建临时二维码,个数不限，最大有效期30天(2592000秒)
 * @param sceneId 场景 ID 值，非0
 * @param sceneStr 场景 string 值
 */
func CreateTemporaryQrCode(sceneId int, sceneStr string, expireSeconds int) (string, error) {
	if expireSeconds <= 0 || expireSeconds > 2592000 {
		return "", util.ErrorParam
	}
	param := qrCodePostInfo{
		Expire: expireSeconds,
	}
	//优先判断字符串场景值
	if len(sceneStr) != 0 {
		if len(sceneStr) > 64 {
			return "", util.ErrorParam
		} else {
			action := actionInfo{
				Scene: sceneInfo{
					Str: sceneStr,
				},
			}
			param.Name = "QR_STR_SCENE"
			param.Info = action
		}
	} else {
		if sceneId <= 0 {
			return "", util.ErrorParam
		} else {
			action := actionInfo{
				Scene: sceneInfo{
					Id: sceneId,
				},
			}
			param.Name = "QR_SCENE"
			param.Info = action
		}
	}
	ticket, err := createQrCode(param)
	if err != nil {
		return "", err
	}
	return url.QueryEscape(ticket), err
}

//return openId
func CheckScanQrCode(ticket string) (string, error) {
	if len(ticket) == 0 {
		return "", util.ErrorParam
	}
	ticket, err := url.QueryUnescape(ticket)
	if err != nil {
		return "", util.ErrorParam
	}
	conn := dao.RedisPool.Get()
	defer conn.Close()
	id, err := conn.Do("get", util.WechatTicketPrefix+ticket)
	if err != nil || id == nil {
		return "", util.ErrorQrcodeNotScan
	}
	return string(id.([]byte)), nil
}

func createQrCode(param qrCodePostInfo) (ticket string, err error) {
	data, _ := json.Marshal(param)
	req := util.Request{
		Method: util.POST,
		Url:    "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=" + model.GetAccessToken(),
		Type:   util.JSON,
		Params: data,
	}
	res := util.Call(req)
	if res.Status != util.StatusOK {
		err = errors.New("status is :" + strconv.Itoa(res.Status))
		return
	} else {
		var result interface{}
		if nil == json.Unmarshal([]byte(res.Body), &result) {
			data := result.(map[string]interface{})
			if t, ok := data["ticket"]; ok {
				ticket = t.(string)
			}
			if len(ticket) == 0 {
				err = util.ErrorHttpResponse
				return
			}
			return
		}
	}
	err = util.ErrorHttpResponse
	return
}
