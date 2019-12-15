package message

import (
	"encoding/xml"
	"time"
	"wechat/util"
)

//微信推送消息
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

type Media struct {
	MediaId string `xml: "MediaId"`
	//仅视频消息需要下面两项
	Title       string `xml: "Title"`
	Description string `xml: "Description"`
}

//图片回复消息
type ImageMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml: "ToUserName"`
	FromUserName string   `xml: "FromUserName"`
	CreateTime   int64    `xml: "CreateTime"`
	MsgType      string   `xml: "MsgType"`
	Image        Media    `xml: "Image"`
}

//音频回复消息
type VoiceMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml: "ToUserName"`
	FromUserName string   `xml: "FromUserName"`
	CreateTime   int64    `xml: "CreateTime"`
	MsgType      string   `xml: "MsgType"`
	Voice        Media    `xml: "Voice"`
}

//视频回复消息
type VideoMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml: "ToUserName"`
	FromUserName string   `xml: "FromUserName"`
	CreateTime   int64    `xml: "CreateTime"`
	MsgType      string   `xml: "MsgType"`
	Video        Media    `xml: "Video"`
}

//文本回复消息
type TextMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml: "ToUserName"`
	FromUserName string   `xml: "FromUserName"`
	CreateTime   int64    `xml: "CreateTime"`
	MsgType      string   `xml: "MsgType"`
	Content      string   `xml: "Content"`
}

type Reply interface {
	Reply(Message) (interface{}, error)
}

func (msg Message) Response() (interface{}, error) {
	switch msg.MsgType {
	case util.WechatEvent:
		switch msg.Event {
		case util.WechatEventSubscribe:
			reply := SubscribeReply{}
			return reply.Reply(msg)
		case util.WechatEventUnsubscribe:
			reply := UnsubscribeReply{}
			return reply.Reply(msg)
		case util.WechatEventScan:
			reply := ScanReply{}
			return reply.Reply(msg)
		}
	case util.WechatText:
		reply := TextReply{}
		return reply.Reply(msg)
	}
	return nil, nil
}

//创建文本消息
func newTextMessage(toUser, fromUser, content string) TextMessage {
	return TextMessage{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}
}

//创建图片消息
func newImageMessage(toUser, fromUser, mediaId string) ImageMessage {
	media := Media{MediaId: mediaId}
	return ImageMessage{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "image",
		Image:        media,
	}
}

//创建音频消息
func newVoiceMessage(toUser, fromUser, mediaId string) VoiceMessage {
	media := Media{MediaId: mediaId}
	return VoiceMessage{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "voice",
		Voice:        media,
	}
}

//创建视频消息
func newVideoMessage(toUser, fromUser, mediaId, title, description string) VideoMessage {
	media := Media{mediaId, title, description}
	return VideoMessage{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "video",
		Video:        media,
	}
}
