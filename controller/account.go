package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"wechat/model/officialAccount/user"
	"wechat/util"
)

func addAccount(c *gin.Context) {
	role, _ := c.Get("role")
	if role.(int) != util.WechatPlatformManager {
		c.JSON(util.StatusOK, genResult(util.ERROR_UNAUTH))
		return
	}

	name := c.PostForm("name")
	phone := c.PostForm("phone")
	status := c.PostForm("status")
	cstatus := 0
	if status == "1" {
		cstatus = 1
	}

	part := c.PostForm("role")
	var crole int = 1
	if len(part) == 0 {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	if part == "1" {
		crole = 1
	} else if part == "2" {
		crole = 2
	} else {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}

	if len(name) != 0 && len(name) > 10 {
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
	account, err := user.AddAccount(name, mobile, crole, cstatus)
	if err != nil {
		result := genResult(util.ERROR_PARAM)
		result["errmsg"] = err.Error()
		c.JSON(util.StatusOK, result)
		return
	}
	result := genResult(util.OK)
	result["account"] = account
	c.JSON(util.StatusOK, result)
}

func accountList(c *gin.Context) {
	status := c.PostForm("status")
	cstatus := 1
	if len(status) == 0 {
		cstatus = -1
	} else if status == "1" {
		cstatus = 1
	} else {
		cstatus = 0
	}

	role := c.PostForm("role")
	crole := 1
	if len(role) == 0 {
		crole = -1
	} else if role == "1" {
		crole = 1
	} else {
		crole = 2
	}

	var pageNo int = 1
	no, err := strconv.Atoi(c.PostForm("pageNo"))
	if err == nil && no > 0 {
		pageNo = no
	}

	var pageSize int = 10
	size, err := strconv.Atoi(c.PostForm("pageSize"))
	if err == nil && size > 0 {
		pageSize = size
	}

	name := c.PostForm("name")
	phone := c.PostForm("phone")
	var mobile int64 = 0
	if len(phone) > 0 {
		if len(phone) != 11 {
			c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
			return
		} else {
			cmobile, err := strconv.ParseInt(phone, 10, 64)
			if err != nil {
				c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
				return
			}
			mobile = cmobile
		}
	}

	result := genResult(util.OK)
	total, accounts, err := user.GetAccountList(name, mobile, crole, cstatus, pageNo, pageSize)
	if err != nil {
		result := genResult(util.ERROR_SYSTEM)
		result["errmsg"] = err.Error()
		c.JSON(util.StatusOK, result)
		return
	}
	result["total"] = total
	result["pageNo"] = pageNo
	result["pageSize"] = pageSize
	result["accounts"] = accounts
	c.JSON(util.StatusOK, result)
}

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

func updateAccountById(c *gin.Context) {
	role, _ := c.Get("role")
	if role.(int) != util.WechatPlatformManager {
		c.JSON(util.StatusOK, genResult(util.ERROR_UNAUTH))
		return
	}

	id := c.PostForm("id")
	if len(id) == 0 {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	uid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}

	name := c.PostForm("name")
	phone := c.PostForm("phone")
	status := c.PostForm("status")
	part := c.PostForm("role")
	if len(name) == 0 && len(phone) == 0 && len(status) == 0 && len(part) == 0 {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	if len(name) != 0 && len(name) > 10 {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}

	var mobile int64 = 0
	if len(phone) != 0 {
		if len(phone) != 11 {
			c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
			return
		}
		cmobile, err := strconv.ParseInt(phone, 10, 64)
		if err != nil {
			c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
			return
		}
		mobile = cmobile
	}
	cstatus := 0
	if len(status) != 0 {
		if status != "0" && status != "1" {
			c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
			return
		}
		cs, err := strconv.Atoi(status)
		if err != nil {
			c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
			return
		}
		cstatus = cs
	}

	crole := 1
	if len(part) != 0 {
		if part != "1" && part != "2" {
			c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
			return
		}
		cr, err := strconv.Atoi(part)
		if err != nil {
			c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
			return
		}
		crole = cr
	}

	account, err := user.UpdateAccountById(uid, name, mobile, crole, cstatus)
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
