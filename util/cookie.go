package util

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func MakeCookie(openId string) (string, error) {
	data := fmt.Sprintf("%d+%s", time.Now().Unix(), openId)
	src, err := AesEncrypt(data)
	if err != nil {
		return "", err
	} else {
		return url.QueryEscape(base64.StdEncoding.EncodeToString(src)), nil
	}
}

func ParseCookie(cookie string) (string, error) {
	if len(cookie) == 0 {
		return "", ErrorParam
	}
	udata, err := url.QueryUnescape(cookie)
	if err != nil {
		return "", err
	}
	bdata, err := base64.StdEncoding.DecodeString(udata)
	if err != nil {
		return "", err
	}
	src, err := AesDecrypt(bdata)
	if err != nil {
		return "", err
	}
	data := strings.Split(src, "+")
	if len(data) != 2 {
		return "", ErrorParam
	}
	openId := data[1]
	if len(openId) == 0 {
		return "", ErrorParam
	}
	return openId, nil
}
