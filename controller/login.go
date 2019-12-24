package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"wechat/model/officialAccount/qrcode"
	"wechat/util"
)

func loginTicket(c *gin.Context) {
	ticket, err := qrcode.CreateTemporaryQrCode(0, util.WechatPlatformLogin, util.WechatTicketExpire)
	if err != nil {
		c.JSON(util.StatusOK, genResult(util.ERROR_SYSTEM))
		return
	}
	result := genResult(util.OK)
	result["ticket"] = ticket
	c.JSON(util.StatusOK, result)
}

func checkScanQrcode(c *gin.Context) {
	ticket := c.PostForm("ticket")
	if len(ticket) == 0 {
		c.JSON(util.StatusOK, genResult(util.ERROR_PARAM))
		return
	}
	openId, err := qrcode.CheckScanQrCode(ticket)
	if err != nil {
		c.JSON(util.StatusOK, genResult(util.ERROR_QRCODE_NOT_SCAN))
		return
	}
	if len(openId) == 0 {
		c.JSON(util.StatusOK, genResult(util.ERROR_SYSTEM))
		return
	}
	cookie, err := util.MakeCookie(openId)
	if err != nil {
		util.Logger.Warn("make cookie fail", zap.String("openId", openId), zap.String("error", err.Error()))
		c.JSON(util.StatusOK, genResult(util.ERROR_SYSTEM))
		return
	}
	result := genResult(util.OK)
	result["cookie"] = cookie
	c.JSON(util.StatusOK, result)
}
