package user

import ()

type Wxer struct {
	Id            int64  `json:"id" db:"id"`
	OpenId        string `json:"open_id" db:"open_id"`
	UnionId       string `json:"union_id" db:"union_id"`
	NickName      string `json:"nick_name" db:"nick_name"`
	HeadImgUrl    string `json:"head_img_url" db:"head_img_url"`
	Sex           int    `json:"sex" db:"sex"`
	City          string `json:"city" db:"city"`
	Province      string `json:"province" db:"province"`
	Country       string `json:"country" db:"country"`
	Language      string `json:"language" db:"language"`
	Uid           int64  `json:"uid" db:"uid"`
	Subscribe     int    `json:"subscribe" db:"subscribe"`
	SubscribeTime int    `json:"subscribe_time" db:"subscribe_time"`
	CreateTime    string `json:"create_time" db:"create_time"`
	UpdateTime    string `json:"update_time" db:"update_time"`
}
