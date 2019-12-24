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

func AddAccount(name string, phone int64, role, status int) (*Account, error) {
	account := new(Account)
	dao.DB.Get(account, "SELECT * FROM account WHERE phone=?", phone)
	if account.Phone != 0 {
		util.Logger.Error("AddAccount fail, phone exist", zap.Int64("phone", phone))
		return nil, util.ErrorPhoneExist
	}
	account.Name = name
	account.Phone = phone
	account.Role = role
	account.Status = status
	_, err := dao.DB.NamedExec("INSERT INTO account(name,phone,role,status) VALUES(:name, :phone, :role, :status)", account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func GetAccountList(name string, phone int64, role, status, pageNo, pageSize int) (int, []Account, error) {
	//获取总数
	var total int
	sql := " FROM account WHERE 1=1"
	if len(name) > 0 {
		sql += " AND name like '%" + name + "%'"
	}
	if phone != 0 {
		sql += fmt.Sprintf(" AND phone = %d", phone)
	}
	if status >= 0 {
		sql += fmt.Sprintf(" AND status = %d", status)
	}
	if role >= 0 {
		sql += fmt.Sprintf(" AND role = %d", role)
	}
	err := dao.DB.Get(&total, "SELECT count(*)"+sql)
	if err != nil {
		util.Logger.Error("GetAccountList total fail", zap.String("error", err.Error()))
		return 0, nil, err
	}

	//分页查询
	start := (pageNo - 1) * pageSize
	if start > total {
		util.Logger.Error("GetAccountList param error", zap.Int("pageNo", pageNo), zap.Int("pageSize", pageSize), zap.Int("total", total))
		return 0, nil, util.ErrorParam
	}
	accounts := []Account{}
	sql = "SELECT * " + sql + fmt.Sprintf(" limit %d, %d", start, pageSize)
	err = dao.DB.Select(&accounts, sql)
	if err != nil {
		util.Logger.Error("GetAccountList fail", zap.Int("start", start), zap.Int("pageSize", pageSize), zap.String("error", err.Error()))
		return 0, nil, err
	}
	return total, accounts, nil
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
	err := dao.DB.Get(account, "SELECT * FROM account WHERE phone=? and status = 1 limit 1", phone)
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
	if len(account.OpenId) != 0 {
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
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", model.GetAccessToken(util.OfficialAccount), openId)
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

func UpdateAccountById(id int64, name string, phone int64, role, status int) (*Account, error) {
	account := new(Account)
	err := dao.DB.Get(account, "SELECT * FROM account WHERE id=?", id)
	if err != nil {
		util.Logger.Error("UpdateAccountById fail", zap.Int64("id", id), zap.String("error", err.Error()))
		return nil, err
	}
	if len(name) != 0 {
		account.Name = name
	}
	if phone != 0 {
		account.Phone = phone
	}
	if status != -1 {
		account.Status = status
	}
	if role != -1 {
		account.Role = role
	}
	//避免误操作
	if account.Status != 1 || account.Role != 1 {
		var total int
		dao.DB.Get(&total, "SELECT COUNT(*) FORM account WHERE role = 1 and status = 1")
		if total <= 1 {
			return nil, util.ErrorUpdateForbidden
		}
	}
	_, err = dao.DB.NamedExec("UPDATE account SET name=:name, phone=:phone, role =:role, status=:status WHERE id=:id", account)
	if err != nil {
		util.Logger.Error("UpdateAccountById fail", zap.String("account", account.String()), zap.String("error", err.Error()))
		return nil, err
	}
	return account, nil
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
