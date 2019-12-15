package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"wechat/controller"
	"wechat/util"
)

var server *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	server = gin.New()
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
