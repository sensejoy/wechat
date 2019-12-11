package message

import (
	"encoding/xml"
	"time"
)

type TextMessage struct {
	ToUserName   string `xml: "ToUserName"`
	FromUserName string `xml: "FromUserName"`
	CreateTime   int64  `xml: "CreateTime"`
	Content      string `xml: "Content"`
	MsgId        string `xml: "MsgId"`
}

type TextResponse struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml: "ToUserName"`
	FromUserName string   `xml: "FromUserName"`
	CreateTime   int64    `xml: "CreateTime"`
	MsgType      string   `xml: "MsgType"`
	Content      string   `xml: "Content"`
}

func (text TextMessage) Reply() (interface{}, error) {
	//TODO 关键字回复
	res := TextResponse{
		ToUserName:   text.FromUserName,
		FromUserName: text.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      "您好，您的消息:" + text.Content + "已收到，请等待客服回复",
	}
	return res, nil
}
