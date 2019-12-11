package main

import (
	"fmt"
	"testing"
	"wechat/model/message"
)

func TestSendCustomMessage(t *testing.T) {
	msg := map[string]interface{}{
		"touser":  "o9Pozv1rfv_298vRr3kh1IOC5bhA",
		"msgtype": "text",
		"text": map[string]string{
			"content": "客服消息测试：hello",
		},
	}
	err := message.SendCustomMessage(msg)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("发送客服消息成功")
	}
}

func TestSendTemplateMessage(t *testing.T) {
	openId := "o9Pozv1rfv_298vRr3kh1IOC5bhA"
	templateId := "bNU4Aq4h-n8FT-RnH6ndxsWy_mw3-_hx2xIl55YTTE4"
	params := map[string]interface{}{
		"name": map[string]string{
			"value": "宋西军",
			"color": "#173177",
		},
	}
	url := ""
	msgId, err := message.SendTemplateMessage(openId, templateId, url, params)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("发送成功，msgId:", msgId)
	}
}
