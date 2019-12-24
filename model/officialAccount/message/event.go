package message

/**
 * desc 处理微信事件消息
 */

import (
	"go.uber.org/zap"
	"wechat/model/dao"
	"wechat/util"
)

type SubscribeReply struct{}
type UnsubscribeReply struct{}
type ScanReply struct{}

//用户关注
func (reply SubscribeReply) Reply(message Message) (interface{}, error) {
	openId := message.FromUserName
	appId := message.ToUserName
	ticket := message.Ticket
	util.Logger.Info("user subscribe:", zap.String("openId", openId), zap.String("appId", appId))
	//如果是登录后台扫码
	if message.EventKey == util.WechatPlatformLogin {
		return updateLogin(openId, appId, ticket), nil
	}
	//TODO:关注后的自动回复
	return nil, nil
}

//用户取消关注
func (reply UnsubscribeReply) Reply(message Message) (interface{}, error) {
	openId := message.FromUserName
	appId := message.ToUserName
	util.Logger.Info("user unsubscribe:", zap.String("openId", openId), zap.String("appId", appId))
	return nil, nil
}

//用户扫码后处理
func (scan ScanReply) Reply(message Message) (interface{}, error) {
	openId := message.FromUserName
	appId := message.ToUserName
	ticket := message.Ticket
	if message.EventKey == util.WechatPlatformLogin {
		return updateLogin(openId, appId, ticket), nil
	}
	return nil, nil
}

func updateLogin(openId, appId, ticket string) interface{} {
	if len(ticket) == 0 {
		return newTextMessage(openId, appId, "系统异常，请刷新后台页面重新扫码登录")
	}
	//写入缓存中
	conn := dao.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("setex", util.WechatTicketPrefix+ticket, util.WechatTicketExpire, openId)
	if err != nil {
		return newTextMessage(openId, appId, "系统异常，请刷新后台页面重新扫码登录")
	}
	return newTextMessage(openId, appId, "后台扫码成功，请开始您的操作")
}
