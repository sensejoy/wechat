package message

import (
	"encoding/json"
	"encoding/xml"
	"wechat/util"
)

type Message struct {
	//通用信息
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml: "ToUserName"`
	FromUserName string   `xml: "FromUserName"`
	CreateTime   int64    `xml: "CreateTime"`
	//用户消息
	MsgType      string `xml: "MsgType"`
	MsgId        string `xml: "MsgId"`
	Content      string `xml: "Content"`
	PicUrl       string `xml: "PicUrl"`
	MediaId      string `xml: "MediaId"`
	Format       string `xml: "Format"`
	Recognition  string `xml: "Recognition"`
	ThumbMediaId string `xml: "ThumbMediaId"`
	Location_X   string `xml: "Location_X"`
	Location_Y   string `xml: "Location_Y"`
	Scale        string `xml: "Scale"`
	Label        string `xml: "Label"`
	Title        string `xml: "Title"`
	Description  string `xml: "Description"`
	Url          string `xml: "Url"`
	MsgID        string `xml: "MsgID"`
	Status       string `xml: "Status"`
	//用户事件
	Event     string `xml: "Event"`
	EventKey  string `xml: "EventKey"`
	Ticket    string `xml: "Ticket"`
	Latitude  string `xml: "Latitude"`
	Longitude string `xml: "Longitude"`
	Precision string `xml: "Precision"`
}

type Reply interface {
	Reply() (interface{}, error)
}

func (msg Message) Copy() Reply {
	data, _ := json.Marshal(msg)
	switch msg.MsgType {
	case util.WechatText:
		text := TextMessage{}
		json.Unmarshal(data, &text)
		return text
	default:
		return nil
	}
}
