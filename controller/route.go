package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sort"
	"wechat/model/message"
	"wechat/util"
)

func InitRoute(server *gin.Engine) {
	server.GET("/", sample)
	server.GET("/wx/response", response)
	server.POST("/wx/response", callback)
}

func sample(c *gin.Context) {
	c.JSON(util.StatusOK, gin.H{
		"errno":  util.OK,
		"errmsg": util.GetMessage(util.OK),
	})
}

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
		util.Logger.Debug("user send:", zap.String("request_id", c.GetHeader("X-REQUEST-ID")), zap.String("content", string(msg.Content)))
		reply := msg.Copy()
		if reply != nil {
			response, err := reply.Reply()
			if err == nil {
				util.Logger.Warn("reply to user error", zap.String("request_id", c.GetHeader("X-REQUEST-ID")), zap.String("error", err.Error()))
				c.XML(util.StatusOK, response)
			} else {
				c.String(util.StatusOK, err.Error())
			}
		} else {
			c.String(util.StatusOK, "success")
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
