package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
	. "wechat/util"
)

var DB *sqlx.DB

func init() {
	user := Conf["mysql"]["user"].(string)
	pass := Conf["mysql"]["password"].(string)
	server := Conf["mysql"]["server"].(string)
	port := Conf["mysql"]["port"].(int)
	database := Conf["mysql"]["database"].(string)
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", user, pass, server, port, database)
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("open mysql failed,", err)
		os.Exit(1)
	}
	db.SetMaxIdleConns(Conf["mysql"]["MaxIdleConns"].(int))
	db.SetMaxOpenConns(Conf["mysql"]["MaxOpenConns"].(int))
	db.SetConnMaxLifetime(time.Duration(Conf["mysql"]["ConnMaxLifetime"].(int)))
	DB = db
}
