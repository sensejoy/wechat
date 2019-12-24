package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sort"
	"wechat/model/officialAccount/message"
	"wechat/util"
)

//微信校验
func response(c *gin.Context) {
	if valid(c) {
		c.String(util.StatusOK, c.Query("echostr"))
	} else {
		c.String(util.StatusBadRequest, "")
	}
}

//微信回调接口
func callback(c *gin.Context) {
	if valid(c) {
		var msg message.Message
		if nil != c.Bind(&msg) {
			c.String(util.StatusBadRequest, "")
		}
		data, _ := json.Marshal(msg)
		util.Logger.Debug("user send:", zap.String("request_id", c.GetHeader("X-REQUEST-ID")), zap.String("message", string(data)))
		response, err := msg.Response()
		if err == nil {
			c.XML(util.StatusOK, response)
		} else {
			util.Logger.Warn("reply to user error", zap.String("request_id", c.GetHeader("X-REQUEST-ID")), zap.String("error", err.Error()))
			c.String(util.StatusOK, "")
		}
	} else {
		c.String(util.StatusBadRequest, "")
	}
}

func valid(c *gin.Context) bool {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	str := []string{util.Conf["wechat"]["token"].(string), timestamp, nonce}
	sort.Strings(str)
	src := str[0] + str[1] + str[2]
	check := util.Sha1(src)
	return check == signature
}
