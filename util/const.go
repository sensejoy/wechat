package util

import (
	"errors"
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
	ErrorApp        = errors.New("程序未运行")
	ErrorSystem     = errors.New("系统异常")
	ErrorParam      = errors.New("参数异常")
	ErrorParseParam = errors.New("参数解析异常")
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

const App = "wechat"

func GetMessage(errno int) string {
	switch errno {
	case ERROR_LOGIN:
		return "未登录, 请登录后重试"
	case ERROR_USER:
		return "用户异常, 请登录后重试"
	case ERROR_UNKNOWN_USER:
		return "用户不存在, 请登录后重试"
	case ERROR_UNAUTH:
		return "对不起, 您没有权限"
	case ERROR_PARAM:
		return "请求参数异常"
	case ERROR_SYSTEM:
		return "服务器开小差啦，请重试"
	default:
		return ""
	}
}
