package manager

import (
	"encoding/xml"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"manager/models"
	"manager/models/dbs"
	"strings"
	"time"
)

type Crawlers struct{}

var TimeFormat map[string]string = map[string]string{"csdn": "2006/01/02 15:04:05",
	"wordpress": "Mon, 2 Jan 2006 15:04:05 -0700"}

func (this Crawlers) GetCsdn(url string, uid int64) {
	request := gorequest.New()
	defer func() {
		if err := recover(); err != nil {
			log.Println("err：", err)
		}
	}()
	resp, body, errs := request.Get(url).End()
	if errs != nil {
		log.Println(errs)
	}
	var content models.Csdn
	if resp.Status == "200 OK" {
		xml.Unmarshal([]byte(body), &content)
		blogContent := content.Channel.Item
		for i, blog := range blogContent {
			log.Println(i, blog.Title, blog.PubDate)
			tmpBlog := dbs.Blogs{}.GetBlogByBlogArticleLink(blog.Link)
			if tmpBlog.Id != 0 {
				dbs.Blogs{}.UpdateBlog(tmpBlog.Id, blog.Title, blog.Author, this.stringTime2Timestamp(blog.PubDate, "csdn"), blog.Description, nil, nil)
			} else {
				dbs.Blogs{}.InsertBlog(uid, blog.Link, blog.Title, blog.Author, this.stringTime2Timestamp(blog.PubDate, "csdn"), blog.Description, nil, nil)
			}
		}
	}
}

func (this Crawlers) GetWordpress(url string, uid int64) {
	request := gorequest.New()
	defer func() {
		if err := recover(); err != nil {
			log.Println("err：", err)
		}
	}()
	resp, body, errs := request.Get(url).End()
	if errs != nil {
		log.Println(errs)
	}
	var content models.Wordpress
	if resp.Status == "200 OK" {
		xml.Unmarshal([]byte(body), &content)
		blogContent := content.Channel.Item
		for i, blog := range blogContent {
			log.Println(i, blog.Title, blog.PubDate)
			tmpBlog := dbs.Blogs{}.GetBlogByBlogArticleLink(blog.Link)
			if tmpBlog.Id != 0 {
				dbs.Blogs{}.UpdateBlog(tmpBlog.Id, blog.Title, blog.Author, this.stringTime2Timestamp(blog.PubDate, "wordpress"), blog.Description, blog.Summary, strings.Join(blog.Category, ","))
			} else {
				dbs.Blogs{}.InsertBlog(uid, blog.Link, blog.Title, blog.Author, this.stringTime2Timestamp(blog.PubDate, "wordpress"), blog.Description, blog.Summary, strings.Join(blog.Category, ","))
			}
		}
	}
}

func (this Crawlers) stringTime2Timestamp(timer, typer string) int64 {
	format, _ := TimeFormat[typer]
	loc, _ := time.LoadLocation("Local")
	tm, _ := time.ParseInLocation(format, timer, loc)
	fmt.Println(tm.String())
	//TODO 由于历史原因 时间需要 加上 1901 年， 然后转换成毫秒计的时间戳
	return (tm.Unix() + 59958144000) * 1000
}
