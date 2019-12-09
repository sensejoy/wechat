package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"wechat/util"
)

var command string

func init() {
	timestamp := time.Now().Format("2006-01-02-15:04:05")
	command = "./wechat" + " &>log/run." + timestamp + " &"
}

func main() {
	if len(os.Args) != 2 {
		help()
	}
	switch os.Args[1] {
	case "start":
		start()
	case "stop":
		stop()
	case "restart":
		restart()
	default:
		help()
	}
}

func help() {
	fmt.Println("Usage:")
	fmt.Println("	./control [start|stop|restart]")
	os.Exit(0)
}

func start() {
	if util.CheckRun() != nil {
		fmt.Println("程序已经启动")
		return
	}
	cmd := exec.Command("/bin/bash", "-c", command)
	if err := cmd.Start(); err != nil {
		fmt.Println("启动失败:", err)
	} else {
		fmt.Println("程序启动成功")
	}
}
func stop() {
	if err := util.KillApp(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("程序成功退出")
	}
}
func restart() {
	if err := util.KillApp(); err != nil {
		fmt.Println(err)
	}
	start()
}
