package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"wechat/model"
	"wechat/model/dao"
	"wechat/util"
)

type Account struct {
	Id     int64  `json:"id" db:"id"`
	OpenId string `json:"open_id" db:"open_id"`
	Name   string `json:"name" db:"name"`
	Phone  int64  `json:"phone" db:"phone"`
	Role   int    `json:"role" db:"role"`
	Avatar string `json:"avatar" db:"avatar"`
	Status int    `json:"status" db:"status"`
	Create string `json:"create_time" db:"create_time"`
	Update string `json:"update_time" db:"update_time"`
}

func GetAccountInfo(openId string) (*Account, error) {
	account := &Account{}
	err := dao.DB.Get(account, "SELECT * FROM account WHERE open_id=? and status = 1", openId)
	if err != nil {
		util.Logger.Error("GetAccountInfo fail", zap.String("openId", openId), zap.String("error", err.Error()))
		return nil, err
	}
	return account, nil
}

func GetAccountInfoByPhone(phone int64) (*Account, error) {
	account := &Account{}
	err := dao.DB.Get(account, "SELECT * FROM account WHERE phone=? and status = 1", phone)
	if err != nil {
		util.Logger.Error("GetAccountInfo fail", zap.Int64("phone", phone), zap.String("error", err.Error()))
		return nil, err
	}
	return account, nil
}

func BindAccount(openId, name string, phone int64) (*Account, error) {
	account, err := GetAccountInfoByPhone(phone)
	if err != nil {
		return nil, err
	}
	if len(account.Name) != 0 {
		return nil, errors.New("该手机号已被" + account.Name + "绑定!")
	}
	avatar, err := GetAccountWechatInfo(openId)
	if err != nil {
		return nil, err
	}
	account.OpenId = openId
	account.Name = name
	account.Avatar = avatar
	_, err = dao.DB.NamedExec("UPDATE account SET name=:name, open_id=:open_id, avatar=:avatar WHERE phone =:phone", account)
	if err != nil {
		util.Logger.Error("BindAccount fail", zap.String("account", account.String()), zap.String("error", err.Error()))
		return nil, err
	}
	return account, nil
}

func GetAccountWechatInfo(openId string) (string, error) {
	if len(openId) == 0 {
		return "", util.ErrorParam
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", model.GetAccessToken(), openId)
	req := util.Request{
		Method: util.GET,
		Url:    url,
	}
	res := util.Call(req)
	if res.Status != util.StatusOK {
		return "", util.ErrorHttpResponse
	}
	var result interface{}
	if nil != json.Unmarshal([]byte(res.Body), &result) {
		return "", util.ErrorWechatResponse
	}
	data := result.(map[string]interface{})
	avatar, ok := data["headimgurl"]
	if !ok {
		msg := ""
		if errmsg, ok := data["errmsg"]; ok {
			msg = errmsg.(string)
		}
		return "", errors.New("微信数据异常:" + msg)
	}
	return avatar.(string), nil
}

func UpdateAccount(openId, name string, phone int64) (*Account, error) {
	account, err := GetAccountInfo(openId)
	if err != nil {
		return nil, err
	}
	account.Name = name
	account.Phone = phone
	_, err = dao.DB.NamedExec("UPDATE account SET name=:name, phone=:phone WHERE open_id=:open_id", account)
	if err != nil {
		util.Logger.Error("UpdateAccount fail", zap.String("account", account.String()), zap.String("error", err.Error()))
		return nil, err
	}
	return account, nil
}

func (account Account) String() string {
	return fmt.Sprintf("[id:%d, open_id:%s, phone:%s, avatar:%s, name:%s, status:%d]", account.Id, account.OpenId, account.Phone, account.Avatar, account.Name, account.Status)
}
