package message

import ()

type TextReply struct{}

//处理用户文本消息
func (reply TextReply) Reply(message Message) (interface{}, error) {
	//TODO 关键字回复
	//return newTextMessage(message.FromUserName, message.ToUserName, "您好，您的消息:"+message.Content+"已收到，请等待客服回复"), nil
	return newImageMessage(message.FromUserName, message.ToUserName, "oJbCdNYIXtEQ-t4ghYC6zju-t9mhFrT6hETtbZfk70CujoRnrfEeg-891hI02FpE"), nil
}
