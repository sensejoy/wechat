package material

/**
 * 临时素材相关:
 * 1、临时素材media_id是可复用的。
 * 2、媒体文件在微信后台保存时间为3天，即3天后media_id失效。
 * 3、上传临时素材的格式、大小限制与公众平台官网一致。
 * 图片（image）: 2M，支持PNG\JPEG\JPG\GIF格式
 * 语音（voice）：2M，播放长度不超过60s，支持AMR\MP3格式
 * 视频（video）：10MB，支持MP4格式
 * 缩略图（thumb）：64KB，支持JPG格式
 */

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"wechat/model"
	"wechat/util"
)

/**
 * 获取临时图片下载到本地
 */
func GetTemporaryImage(mediaId string) (string, error) {
	return "", nil
}

/**
 * desc 上传本地图片, 注意大小2MB以内，支持 PNG\JPEG\JPG\GIF
 */
func UploadImage(filePath, fileName string) (string, error) {
	if err := uploadCheck(filePath, fileName, util.WechatFileImage); err != nil {
		return "", err
	}
	return uploadMedia(filePath, fileName, util.WechatFileImage)
}

/**
 * desc 上传本地音频文件，2MB以内，仅支持 AMR/MP3
 */
func UploadVoice(filePath, fileName string) (string, error) {
	if err := uploadCheck(filePath, fileName, util.WechatFileVoice); err != nil {
		return "", err
	}
	return uploadMedia(filePath, fileName, util.WechatFileVoice)
}

/**
 * desc 上传本地视频文件，10MB以内，仅支持 MP4
 */
func UploadVideo(filePath, fileName string) (string, error) {
	if err := uploadCheck(filePath, fileName, util.WechatFileVideo); err != nil {
		return "", err
	}
	return uploadMedia(filePath, fileName, util.WechatFileVideo)
}

/**
 * desc 上传本地音频文件，64KB以内，仅支持 JPG
 */
func UploadThumb(filePath, fileName string) (string, error) {
	if err := uploadCheck(filePath, fileName, util.WechatFileThumb); err != nil {
		return "", err
	}
	return uploadMedia(filePath, fileName, util.WechatFileThumb)
}

func uploadCheck(filePath, fileName, category string) error {
	filename := filePath + "/" + fileName
	file, err := os.Open(filename)
	if err != nil {
		return util.ErrorFileNotFound
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return err
	}
	ext := strings.ToLower(path.Ext(fileName))

	switch category {
	case util.WechatFileImage:
		switch ext {
		case ".png":
		case ".jpg":
		case ".jpeg":
		case ".gif":
		default:
			return util.ErrorInvalidType
		}
		if fileinfo.Size() > util.IMAGE_MAX_SIZE {
			return util.ErrorFileTooLarge
		}
	case util.WechatFileVoice:
		switch ext {
		case ".amr":
		case ".mp3":
		default:
			return util.ErrorInvalidType
		}
		if fileinfo.Size() > util.VOICE_MAX_SIZE {
			return util.ErrorFileTooLarge
		}
	case util.WechatFileVideo:
		switch ext {
		case ".mp4":
		default:
			return util.ErrorInvalidType
		}
		if fileinfo.Size() > util.VIDEO_MAX_SIZE {
			return util.ErrorFileTooLarge
		}
	case util.WechatFileThumb:
		switch ext {
		case ".jpg":
		default:
			return util.ErrorInvalidType
		}
		if fileinfo.Size() > util.THUMB_MAX_SIZE {
			return util.ErrorFileTooLarge
		}
	}
	return nil
}

/**
 * desc 临时上传本地文件
 * filePath 文件路径
 * fileName 文件名称
 * category 文件类型，目前微信仅支持image、voice、video、thumb
 * return mediaId string
 * return error
 */
func uploadMedia(filePath, fileName, category string) (string, error) {
	params := map[string]string{
		"filePath":  filePath,
		"fileName":  fileName,
		"nameField": "media",
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s", model.GetAccessToken(), category)
	req := util.Request{
		Method: util.POST,
		Url:    url,
		Type:   util.FORMDATA,
		Params: params,
	}
	res := util.Call(req)
	if res.Status != util.StatusOK {
		return "", errors.New(fmt.Sprintf("status not 200: %d", res.Status))
	}
	var result interface{}
	if nil != json.Unmarshal([]byte(res.Body), &result) {
		return "", util.ErrorHttpResponse
	} else {
		data := result.(map[string]interface{})
		if id, ok := data["media_id"]; ok {
			return id.(string), nil
		} else {
			return "", errors.New(data["errmsg"].(string))
		}
	}
}
