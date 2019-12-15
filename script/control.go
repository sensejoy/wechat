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
	command = util.Dir + "/" + util.App + " &>" + util.Dir + "/run/run." + timestamp + " &"
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
	case "status":
		status()
	default:
		help()
	}
}

func help() {
	fmt.Println("Usage:")
	fmt.Println("	" + util.Dir + "/control [start|stop|restart|status]")
	os.Exit(0)
}

func start() {
	if util.CheckRun() != nil {
		fmt.Println(util.App + "已经启动")
		return
	}
	cmd := exec.Command("/bin/bash", "-c", command)
	if err := cmd.Start(); err != nil {
		fmt.Println(util.App+"启动失败:", err)
	}
	ch := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 2)
		ch <- struct{}{}
	}()
	<-ch
	if util.CheckRun() == nil {
		fmt.Println(util.App + "启动异常，请查看运行日志")
	} else {
		fmt.Println(util.App + "启动成功")
	}
}
func stop() {
	if err := util.KillApp(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(util.App + "成功退出")
	}
}
func restart() {
	if err := util.KillApp(); err != nil {
		fmt.Println(err)
	}
	start()
}
func status() {
	if util.CheckRun() != nil {
		fmt.Println(util.App + "运行中")
	} else {
		fmt.Println(util.App + "未启动")
	}
}
