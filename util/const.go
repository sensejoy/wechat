package util

import (
	"errors"
)

const (
	//消息类型
	WechatText        = "text"
	WechatImage       = "image"
	WechatVoice       = "voice"
	WechatMusic       = "music"
	WechatVideo       = "video"
	WechatShortVideo  = "shortvideo"
	WechatLocation    = "location"
	WechatLink        = "link"
	WechatNews        = "news"
	WechatMpNews      = "mpnews"
	WechatMenu        = "msgmenu"
	WechatCard        = "wxcard"
	WechatMiniProgram = "miniprogrampage"

	//输入状态
	WechatCommandTyping       = "Typing"
	WechatCommandCancelTyping = "CancelTyping"

	//事件类型
	WechatEvent            = "event"
	WechatEventSubscribe   = "subscribe"
	WechatEventUnsubscribe = "unsubscribe"
	WechatEventScan        = "SCAN"
	WechatEventLocation    = "LOCATION"
	WechatEventClick       = "CLICK"
	WechatEventView        = "VIEW"

	//登录后台事件值
	WechatPlatformLogin = "wechat_platform_login"

	//文件类型
	WechatFileImage = "image"
	WechatFileVoice = "voice"
	WechatFileVideo = "video"
	WechatFileThumb = "thumb"

	WechatTicketPrefix = "wechat:qrcode:ticket:"
	WechatTicketExpire = 3600 //1小时有效

	WechatToken = "wechat_token"

	WechatPlatformManager  = 1
	WechatPlatformOperator = 2

	OfficialAccount = 1
	MiniProgram     = 2
)

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusMethodNotAllowed    = 405
	StatusInternalServerError = 500
	StatusBadGateway          = 502
	StatusServiceUnavailable  = 503
	StatusGatewayTimeout      = 504
)

var (
	ErrorApp             = errors.New("程序未运行")
	ErrorSystem          = errors.New("系统异常")
	ErrorParam           = errors.New("参数异常")
	ErrorParseParam      = errors.New("参数解析异常")
	ErrorInvalidType     = errors.New("文件格式不正确")
	ErrorFileNotFound    = errors.New("文件不存在")
	ErrorFileTooLarge    = errors.New("文件尺寸太大")
	ErrorHttpResponse    = errors.New("错误的响应数据")
	ErrorQrcodeNotScan   = errors.New("二维码未扫描")
	ErrorWechatResponse  = errors.New("微信请求异常请重试")
	ErrorPhoneExist      = errors.New("手机号已经存在")
	ErrorUpdateForbidden = errors.New("因管理员仅余一名，当前操作被禁止")
)

const (
	OK = iota
	ERROR_LOGIN
	ERROR_USER
	ERROR_UNKNOWN_USER
	ERROR_UNAUTH
	ERROR_PARAM
	ERROR_SYSTEM
)

const (
	ERROR_QRCODE_NOT_SCAN = 200 + iota
)

const App = "wechat"

const (
	IMAGE_MAX_SIZE = 1024 * 1024 * 2
	VOICE_MAX_SIZE = 1024 * 1024 * 2
	VIDEO_MAX_SIZE = 1024 * 1024 * 10
	THUMB_MAX_SIZE = 1024 * 64
)

func GetMessage(errno int) string {
	switch errno {
	case ERROR_LOGIN:
		return "未登录, 请登录后重试"
	case ERROR_USER:
		return "用户异常, 请登录后重试"
	case ERROR_UNKNOWN_USER:
		return "用户不存在, 请绑定后重试"
	case ERROR_UNAUTH:
		return "对不起, 您没有权限"
	case ERROR_PARAM:
		return "请求参数异常"
	case ERROR_SYSTEM:
		return "服务器开小差啦，请重试"
	case ERROR_QRCODE_NOT_SCAN:
		return "二维码未扫描"
	default:
		return ""
	}
}
