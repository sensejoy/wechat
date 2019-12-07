package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sort"
	"time"
	"wechat/util"
)

var server *gin.Engine

func init() {
	server = gin.New()
	server.Use(requestInit(), gin.Recovery(), log())
	initRoute()
}

func initRoute() {
	server.GET("/", sample)
	server.GET("/wx/response", response)
}

func main() {
	if err := server.Run(":8808"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
	str := []string{util.Wechat_Token, timestamp, nonce}
	sort.Strings(str)
	src := str[0] + str[1] + str[2]
	check := util.Sha1(src)
	if check == signature {
		c.String(util.StatusOK, echostr)
	} else {
		c.String(util.StatusBadRequest, "bad params")
	}
}

func requestInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("request_id", c.GetHeader("X_REQUEST_ID"))
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
		requestId, _ := c.Get("request_id")
		util.Logger.Info("http request", zap.String("request_id", requestId.(string)), zap.Int("status", c.Writer.Status()), zap.String("client_ip", c.ClientIP()), zap.String("cost", fmt.Sprintf("%dms", cost)), zap.String("request", formatString(r)))
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
