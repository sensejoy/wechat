package message

import (
	"encoding/json"
	"errors"
	"strconv"
	"wechat/model"
	"wechat/util"
)

/**
 * desc 发送客服消息
 *
 */
func SendCustomMessage(msg map[string]interface{}) error {
	data, _ := json.Marshal(msg)
	req := util.Request{
		Method: util.POST,
		Url:    "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + model.GetAccessToken(),
		Type:   util.JSON,
		Params: data,
	}
	res := util.Call(req)
	if res.Status != util.StatusOK {
		return errors.New("status is :" + strconv.Itoa(res.Status))
	} else {
		var result interface{}
		if nil == json.Unmarshal([]byte(res.Body), &result) {
			data := result.(map[string]interface{})
			if code, ok := data["errcode"]; ok {
				if code.(float64) == 0 {
					return nil
				} else {
					return errors.New(data["errmsg"].(string))
				}
			}
		}
	}
	return errors.New("invalid response")
}

/**
 * desc 发送模板消息
 * openId 发送用户的open-id
 * templateId 模板ID
 * url 模板消息跳转连接，可为空
 * params 模板对应的参数信息
 * return msgid 发送成功微信记录的消息id, 后续微信会推送此消息id触达用户的状态
 * return error 如果发送失败返回错误
 */
func SendTemplateMessage(openId, templateId, url string, param map[string]interface{}) (msgid int64, err error) {
	params := map[string]interface{}{
		"touser":      openId,
		"template_id": templateId,
		"url":         url,
		"data":        param,
	}
	data, _ := json.Marshal(params)
	req := util.Request{
		Method: util.POST,
		Url:    "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + model.GetAccessToken(),
		Type:   util.JSON,
		Params: data,
	}
	res := util.Call(req)
	if res.Status != util.StatusOK {
		err = errors.New("status is :" + strconv.Itoa(res.Status))
	} else {
		var result interface{}
		if nil != json.Unmarshal([]byte(res.Body), &result) {
			err = util.ErrorHttpResponse
		} else {
			data := result.(map[string]interface{})
			if id, ok := data["msgid"]; ok {
				msgid = int64(id.(float64))
			} else {
				err = errors.New(data["errmsg"].(string))
			}
		}
	}
	return
}
