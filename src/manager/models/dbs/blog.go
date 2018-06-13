package dbs

import (
	"fmt"
	"common"
)

type Blog struct {
	Id int64 `db:"id" json:"id"`
	Uid int64 `db:"uid" json:"uid"`
	BlogArticleLink string `db:"BlogArticleLink" json:"blogArticleLink"`
	Title string `db:"Title" json:"title"`
	Author string `db:"Author" json:"author"`
	PubDate int64 `db:"PubDate" json:"pubDate"`
	ArticleDetail string `db:"ArticleDetail" json:"articleDetail"`
	Summary interface{} `db:"Summary" json:"summary"`
	Category interface{} `db:"Category" json:"catedory"`
}

type Blogs struct {}

func (this Blogs) GetBlogByBlogArticleLink(link string) *Blog {
	var row Blog
	sql := fmt.Sprintf("select %s from %s where BlogArticleLink=?", this.columns(), this.table())
	common.DB().Get(&row, sql, link)
	return &row
}

func (this Blogs) InsertBlog(uid  int64, link, title, author string, pubDate int64, detail string, summary, category interface{}) {
	sql := fmt.Sprintf("insert into %s (uid, BlogArticleLink, Title, Author, PubDate, ArticleDetail, Summary, Category) values(?, ?, ?, ?, ?, ?, ?, ?)", this.table())
	common.DB().Insert(sql, uid, link, title, author, pubDate, detail, summary, category)
}

func (this Blogs) UpdateBlog(id int64, title, author string, pubDate int64, detail string, summary, category interface{}) {
	sql := fmt.Sprintf("update %s set  Title=?, Author=?, PubDate=?, ArticleDetail=?, Summary=?, Category=? where id=?", this.table())
	common.DB().Update(sql, title, author, pubDate, detail, summary, category, id)
}

// 获取表名
func (this Blogs) table() string {
	return "T_blog"
}

//获取数据表的列

func (this Blogs) columns() string {
	return `id, uid, BlogArticleLink, Title, Author, PubDate, ArticleDetail, Summary, Category`
}
