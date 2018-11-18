package models

import (
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"regexp"
	"strings"
	"testing"
)

func TestCsdn_UploadQiniu(t *testing.T) {
	url := "https://img-blog.csdn.net/201703121218301903"
	splitData := strings.Split(url, "/")
	key := "default"
	if len(splitData) >= 4 {
		key = splitData[3]
	} else {
		return
	}
	//cmdUrl := fmt.Sprintf("curl -H \"Referer:http://blog.csdn.net\" %s -o /tmp/%s.csdn.tmp", url, key)
	//cmd := exec.Command("bash", "-c", cmdUrl)
	//_, err := cmd.CombinedOutput()
	//if err != nil {
	//	log.Printf("cmdUrl:%s, err:%v", cmdUrl, err)
	//}
	//localFile := fmt.Sprintf("/tmp/%s.csdn.tmp", key)

	bucket := "blog"

	mac := qbox.NewMac("BIW84ekVdcZOkDLKJJytWAxlb37RFxlrsQn0SsTA", "K2JkK1JXNfslhu6Czi_PDlqYSbrMM69mH0ohIgYP")
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	fileInfo, sErr := bucketManager.Stat(bucket, key)
	if sErr != nil {
		fmt.Println("err:", sErr)
		return
	}
	fmt.Println(fileInfo.String())
}

func TestCsdn_ReplaceImgUrlToQiniuCdnUrl(t *testing.T) {
	con := `<p><img src="https://img-blog.csdn.net/20180731092048866?watermark/2/text/aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3podXhpbnF1YW42MQ==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70" alt="这里写图片描述" title=""></p>`
	reg := regexp.MustCompile(`https://img-blog\.csdn\.net/.*?\"`)
	result := reg.Find([]byte(con))
	fmt.Println(string(result))
	//imgUrl := "https://img-blog.csdnimg.cn/20181108214953776.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3dlaXhpbl8zNjg4ODU3Nw==,size_16,color_FFFFFF,t_70"
}
