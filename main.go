package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
	"wechat/controller"
	"wechat/util"
)

var server *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	server = gin.New()
	server.Use(requestInit(), gin.Recovery(), log())
	controller.InitRoute(server)
}

func main() {
	if err := util.WritePid(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := server.Run(fmt.Sprintf(":%d", util.Conf["server"]["port"].(int))); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
