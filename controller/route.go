package controller

import (
	"github.com/gin-gonic/gin"
	"sort"
	"wechat/util"
)

func InitRoute(server *gin.Engine) {
	server.GET("/", sample)
	server.GET("/wx/response", response)
}

func sample(c *gin.Context) {
	c.JSON(util.StatusOK, gin.H{
		"errno":  util.OK,
		"errmsg": util.GetMessage(util.OK),
	})
}

func response(c *gin.Context) {
	signature := c.Query("signature")
	echostr := c.Query("echostr")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	str := []string{util.Conf["wechat"]["token"].(string), timestamp, nonce}
	sort.Strings(str)
	src := str[0] + str[1] + str[2]
	check := util.Sha1(src)
	if check == signature {
		c.String(util.StatusOK, echostr)
	} else {
		c.String(util.StatusBadRequest, "bad params")
	}
}
