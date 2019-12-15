package main

import (
	"fmt"
	"testing"
	"wechat/model/qrcode"
	"wechat/util"
)

func TestCreateTemporaryQrcode(t *testing.T) {
	url, err := qrcode.CreateTemporaryQrCode(0, util.WechatPlatformLogin, 600)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("创建临时二维码成功，url:", url)
	}
}
