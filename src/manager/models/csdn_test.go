package models

import "testing"

func TestCsdn_UploadQiniu(t *testing.T) {
	Csdn{}.UploadQiniu("https://img-blog.csdn.net/20170312121830903")
}

//func TestCsdn_ReplaceImgUrlToQiniuCdnUrl(t *testing.T) {
//	Csdn{}.ReplaceImgUrlToQiniuCdnUrl()
//}
