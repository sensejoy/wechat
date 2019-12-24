package main

import (
	"encoding/json"
	"fmt"
	"wechat/model"
	"wechat/util"
)

func main() {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", util.Conf["officialAccount"]["appId"], util.Conf["officialAccount"]["appSecret"])
	req := util.Request{
		Method:           util.GET,
		Url:              url,
		ConnectTimeout:   10000, //one second
		ReadWriteTimeout: 10000,
	}
	res := util.Call(req)
	checkAndSet(res, util.OfficialAccount)

	req.Url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", util.Conf["miniProgram"]["appId"], util.Conf["miniProgram"]["appSecret"])
	res = util.Call(req)
	checkAndSet(res, util.MiniProgram)
}

func checkAndSet(res *util.Response, kind int) {
	fmt.Println(res)
	if res.Status != util.StatusOK {
		fmt.Println("get wechat platform token fail")
	} else {
		var result interface{}
		json.Unmarshal([]byte(res.Body), &result)
		data := result.(map[string]interface{})
		if accessToken, ok := data["access_token"]; ok {
			model.SetAccessToken(accessToken.(string), kind)
		} else {
			fmt.Println("set wechat platform token fail")
		}
	}

}
