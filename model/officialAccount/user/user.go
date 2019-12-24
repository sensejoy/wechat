package user

import (
	"wechat/model/dao"
)

//本地用户表
type User struct {
	Id     int64  `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Phone  int64  `json:"phone" db:"phone"`
	Create string `json:"create_time" db:"create_time"`
	Update string `json:"update_time" db:"update_time"`
}

func AddUser(name string, phone int64) (*User, error) {
	user := new(User)
	dao.DB.Get(user, "SELECT * FROM user WHERE phone=? limit 1", phone)
	if user.Phone != 0 {
		return user, nil
	}
	user.Name = name
	user.Phone = phone
	_, err := dao.DB.NamedExec("INSERT INTO user(name,phone) VALUES(:name, :phone)", user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
