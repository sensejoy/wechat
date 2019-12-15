package main

import (
	"fmt"
	"testing"
	"wechat/model/material"
)

func TestUploadImage(t *testing.T) {
	filePath, fileName := "/tmp", "1.jpg"
	mediaId, err := material.UploadImage(filePath, fileName)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("上传成功，media_id:", mediaId)
	}
}
