package dao

import (
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
	. "wechat/util"
)

var RedisPool *redis.Pool

func init() {
	RedisPool = &redis.Pool{
		MaxIdle:         Conf["redis"]["MaxIdle"].(int),
		MaxActive:       Conf["redis"]["MaxActive"].(int),
		IdleTimeout:     time.Duration(Conf["redis"]["IdleTimeout"].(int)) * time.Second,
		MaxConnLifetime: time.Duration(Conf["redis"]["MaxConnLifetime"].(int)) * time.Second,
		Wait:            Conf["redis"]["Wait"].(bool),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", Conf["redis"]["server"].(string)+":"+strconv.Itoa(Conf["redis"]["port"].(int)))
			if err != nil {
				Logger.Error("redis dial failed", zap.String("error", err.Error()))
				return nil, err
			}
			auth := Conf["redis"]["auth"].(string)
			if len(auth) != 0 {
				if _, err := c.Do("AUTH"); err != nil {
					Logger.Error("redis auth failed", zap.String("error", err.Error()))
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			Logger.Error("redis PING failed", zap.String("error", err.Error()))
			return err
		},
	}
}
