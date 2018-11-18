package models

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type Csdn struct {
	Rss     xml.Name `xml: "rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Channel     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Image       Image    `xml:"iamge"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	Language    string   `xml:"language"`
	Generator   string   `xml:"generator"`
	TTL         string   `xml:"ttl"`
	Copyright   string   `xml:"copyright"`
	PubDate     string   `xml:"pubDate"`
	Item        []Item   `xml:"item"`
}

type Image struct {
	Image xml.Name `xml:"image"`
	Link  string   `xml:"link"`
	Url   string   `xml:"url"`
}

type Item struct {
	Item        xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Guid        string   `xml:"guid"`
	Author      string   `xml:"author"`
	PubDate     string   `xml:"pubDate"`
	Description string   `xml:"description"`
	Category    string   `xml:"category"`
}

func (this Csdn) ReplaceImgUrlToQiniuCdnUrl(content string) string {
	imgUrlReg := regexp.MustCompile(`src=\"https://img-blog\.(csdn|csdnimg)\.(net|cn)/.*?\"`)
	for _, urlData := range imgUrlReg.FindAllString(content, -1) {
		url := string(urlData)
		spliteData := strings.Split(url, "?")
		imgLink := strings.TrimPrefix(spliteData[0], "src=\"")
		watermark := strings.TrimSuffix(spliteData[1], "\"")
		fmt.Println(imgLink, watermark)
		this.UploadQiniu(imgLink)
		replaceSource := fmt.Sprintf("%s?%s", imgLink, watermark)
		replaceTarget := fmt.Sprintf("%s%s%s", "http://blog-image.xiyoulinux.org/", strings.Split(imgLink, "/")[3], "?watermark/2/text/aHR0cDovL2Jsb2cueGl5b3VsaW51eC5vcmc=/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70")
		content = strings.Replace(content, replaceSource, replaceTarget, -1)
	}
	return content
}

func (this Csdn) UploadQiniu(url string) {
	splitData := strings.Split(url, "/")
	key := "default"
	if len(splitData) >= 4 {
		key = splitData[3]
	} else {
		return
	}

	bucket := "blog"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	//mac := qbox.NewMac("BIW84ekVdcZOkDLKJJytWAxlb37RFxlrsQn0SsTA", "K2JkK1JXNfslhu6Czi_PDlqYSbrMM69mH0ohIgYP")
	mac := qbox.NewMac(Conf().Qiniu.AccessKey, Conf().Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	bucketManager := storage.NewBucketManager(mac, &cfg)
	_, sErr := bucketManager.Stat(bucket, key)
	if sErr == nil {
		return
	}

	cmdUrl := fmt.Sprintf("curl -H \"Referer:http://blog.csdn.net\" %s -o /tmp/%s.csdn.tmp", url, key)
	cmd := exec.Command("bash", "-c", cmdUrl)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmdUrl:%s, err:%v", cmdUrl, err)
	}
	localFile := fmt.Sprintf("/tmp/%s.csdn.tmp", key)

	resumeUploader := storage.NewResumeUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.RputExtra{}
	err = resumeUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
}
