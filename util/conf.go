package util

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/kardianos/osext"
	"os"
	"path"
)

var Conf map[string]map[string]interface{}
var Dir, Rdir string

func init() {
	Conf = make(map[string]map[string]interface{})
	dir, err := osext.ExecutableFolder()
	if err != nil {
		fmt.Println("get running path fail:", err)
		os.Exit(1)
	}

	file := path.Join(dir, "/conf/wechat.conf")
	if _, err := os.Stat(file); err != nil {
		fmt.Println(dir + "/conf/wechat.conf not found, try /tmp/wechat.conf")
		file = "/tmp/wechat.conf"
		if _, err := os.Stat(file); err != nil {
			fmt.Println("/tmp/wechat.conf not found")
			os.Exit(1)
		}
	}
	conf, err := ini.Load(file)
	if err != nil {
		fmt.Println(dir + "/conf/wechat.conf load failed")
		os.Exit(1)
	}

	Conf["server"] = map[string]interface{}{
		"port": conf.Section("server").Key("port").RangeInt(8008, 1025, 65535),
	}

	Conf["log"] = map[string]interface{}{
		"Filename":   dir + "/log/wechat.log",
		"MaxSize":    conf.Section("log").Key("maxSize").RangeInt(128, 100, 1000),
		"MaxBackups": conf.Section("log").Key("maxBackups").RangeInt(300, 100, 500),
		"MaxAge":     conf.Section("log").Key("maxAge").RangeInt(7, 0, 30),
		"Compress":   conf.Section("log").Key("compress").MustBool(true),
	}

	Conf["redis"] = map[string]interface{}{
		"server":          conf.Section("redis").Key("server").MustString("127.0.0.1"),
		"port":            conf.Section("redis").Key("port").MustInt(6379),
		"auth":            conf.Section("redis").Key("auth").MustString(""),
		"MaxIdle":         conf.Section("redis").Key("maxIdle").RangeInt(5, 0, 100),
		"MaxActive":       conf.Section("redis").Key("maxActive").RangeInt(5, 0, 200),
		"IdleTimeout":     conf.Section("redis").Key("idleTimeout").RangeInt(120, 0, 3600),
		"MaxConnLifetime": conf.Section("redis").Key("maxConnLifetime").RangeInt(0, 0, 3600),
		"Wait":            conf.Section("redis").Key("wait").MustBool(true),
	}
	Conf["mysql"] = map[string]interface{}{
		"server":          conf.Section("mysql").Key("server").MustString("127.0.0.1"),
		"port":            conf.Section("mysql").Key("port").MustInt(3306),
		"user":            conf.Section("mysql").Key("user").MustString(""),
		"password":        conf.Section("mysql").Key("password").MustString(""),
		"database":        conf.Section("mysql").Key("database").MustString(""),
		"MaxIdle":         conf.Section("mysql").Key("password").MustString(""),
		"MaxIdleConns":    conf.Section("mysql").Key("maxIdleConns").RangeInt(5, 0, 100),
		"MaxOpenConns":    conf.Section("mysql").Key("MaxOpenConns").RangeInt(10, 0, 200),
		"ConnMaxLifetime": conf.Section("redis").Key("connMaxLifetime").RangeInt(0, 0, 3600),
	}

	Conf["wechat"] = map[string]interface{}{
		"appId":     conf.Section("wechat").Key("appId").MustString(""),
		"appSecret": conf.Section("wechat").Key("appSecret").MustString(""),
		"token":     conf.Section("wechat").Key("token").MustString(""),
	}
	Dir = dir
	Rdir = Dir + "/run"
}
