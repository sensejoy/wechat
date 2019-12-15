package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"wechat/model/user"
	"wechat/util"
)

func accountInfo(c *gin.Context) {
	account, ok := c.Get("account")
	if !ok {
		c.JSON(util.StatusOK, genResult(util.ERROR_USER))
	} else {
		result := genResult(util.OK)
		result["account"] = account
		c.JSON(util.StatusOK, result)
	}
}

func bindAccount(c *gin.Context) {
	id, _ := c.Get("openId")
	openId := id.(string)
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	if len(name) == 0 || len(name) > 10 {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	if len(phone) != 11 {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	mobile, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	account, err := user.BindAccount(openId, name, mobile)
	if err != nil {
		result := genResult(util.ERROR_SYSTEM)
		result["errmsg"] = err.Error()
		c.JSON(util.StatusOK, result)
		return
	}
	result := genResult(util.OK)
	result["account"] = account
	c.JSON(util.StatusOK, result)
}

func updateAccount(c *gin.Context) {
	id, _ := c.Get("openId")
	openId := id.(string)
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	if len(name) == 0 || len(name) > 10 {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	if len(phone) != 11 {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	mobile, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	account, err := user.UpdateAccount(openId, name, mobile)
	if err != nil {
		result := genResult(util.ERROR_SYSTEM)
		result["errmsg"] = err.Error()
		c.JSON(util.StatusOK, result)
		return
	}
	result := genResult(util.OK)
	result["account"] = account
	c.JSON(util.StatusOK, result)
}
