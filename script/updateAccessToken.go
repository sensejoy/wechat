package main

import (
	"encoding/json"
	"fmt"
	"wechat/model"
	"wechat/util"
)

func main() {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", util.Conf["wechat"]["appId"], util.Conf["wechat"]["appSecret"])
	req := util.Request{
		Method:           util.GET,
		Url:              url,
		ConnectTimeout:   10000, //one second
		ReadWriteTimeout: 10000,
	}
	res := util.Call(req)
	fmt.Println(url, res)
	if res.Status != util.StatusOK {
		fmt.Println("get wechat token fail")
	} else {
		var result interface{}
		json.Unmarshal([]byte(res.Body), &result)
		data := result.(map[string]interface{})
		if accessToken, ok := data["access_token"]; ok {
			model.SetAccessToken(accessToken.(string))
		} else {
			fmt.Println("get wechat token fail")
		}
	}
}
