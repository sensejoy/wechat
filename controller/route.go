package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"wechat/model/officialAccount/user"
	"wechat/util"
)

var authed = map[string]struct{}{
	"/api/pc/accountinfo": struct{}{},
	"/api/pc/bindaccount": struct{}{},
	"/api/pc/updateaccount": struct{}{},
	"/api/pc/accountlist": struct{}{},
	"/api/pc/updateaccountbyid": struct{}{},
	"/api/pc/addaccount": struct{}{},
}

func InitRoute(server *gin.Engine) {
	server.Use(requestInit(), auth(), gin.Recovery(), log())
	server.GET("/", sample)
	server.GET("/wx/response", response)  //微信验证
	server.POST("/wx/response", callback) //微信回调

	server.GET("/api/pc/loginticket", loginTicket)          //后台登录获取二维码 ticket
	server.POST("/api/pc/checkscanqrcode", checkScanQrcode) //校验扫码情况

	server.GET("/api/pc/accountinfo", accountInfo)      //获取账户信息,需要登录
	server.POST("/api/pc/bindaccount", bindAccount)     //绑定账户，需要登录
	server.POST("/api/pc/updateaccount", updateAccount) //更新账户信息，需要登录

	server.POST("/api/pc/accountlist", accountList)             //查看账户列表，需要登录
	server.POST("/api/pc/updateaccountbyid", updateAccountById) //根据账户 id 更新账户信息，需要管理员权限
	server.POST("/api/pc/addaccount", addAccount)               //添加账户，需要管理员权限

}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if _, ok := authed[path]; !ok {
			c.Next()
		} else {
			token, err := c.Cookie(util.WechatToken)
			if err != nil {
				c.JSON(util.StatusOK, genResult(util.ERROR_LOGIN))
				c.Abort()
				return
			}
			openId, err := util.ParseCookie(token)
			if err != nil {
				c.JSON(util.StatusOK, genResult(util.ERROR_USER))
				c.Abort()
				return
			} else {
				if path != "/api/pc/bindaccount" {
					account, err := user.GetAccountInfo(openId)
					if err != nil {
						c.JSON(util.StatusOK, genResult(util.ERROR_UNKNOWN_USER))
						c.Abort()
						return
					}
					c.Set("account", account)
					c.Set("openId", account.OpenId)
					c.Set("role", account.Role)
				} else {
					c.Set("openId", openId)
				}
			}
			c.Next()
		}
	}
}

func requestInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("start", time.Now().UnixNano())
		c.Next()
	}
}

func log() gin.HandlerFunc {
	return func(c *gin.Context) {
		end := time.Now().UnixNano()
		start, _ := c.Get("start")
		cost := (int)((end - start.(int64)) / 1000000)
		r := c.Request
		util.Logger.Info("http request", zap.String("request_id", c.GetHeader("X-REQUEST-ID")), zap.Int("status", c.Writer.Status()), zap.String("client_ip", c.ClientIP()), zap.String("cost", fmt.Sprintf("%dms", cost)), zap.String("request", formatString(r)))
	}
}

func formatString(r *http.Request) string {
	header, _ := json.Marshal(r.Header)
	cookie, _ := json.Marshal(r.Cookies())
	if r.Method != http.MethodGet {
		r.ParseForm()
		params, _ := json.Marshal(r.PostForm)
		return fmt.Sprintf("method[%s] url[%s] header[%s] cookie[%s] params[%s]", r.Method, r.URL, header, cookie, params)
	} else {
		return fmt.Sprintf("method[%s] url[%s] header[%s] cookie[%s]", r.Method, r.URL, header, cookie)
	}
}

func genResult(errno int) gin.H {
	return gin.H{
		"errno":  errno,
		"errmsg": util.GetMessage(errno),
	}
}

func sample(c *gin.Context) {
	c.JSON(util.StatusOK, gin.H{
		"errno":  util.OK,
		"errmsg": util.GetMessage(util.OK),
	})
}
