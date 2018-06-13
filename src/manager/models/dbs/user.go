package dbs

import (
	"fmt"
	"common"
)

type User struct {
	Id int64 `db:"Id" json:"id"`
	Name string `db:"Name" json:"name"`
	QQ string `db:"QQ" json:"qq"`
	BlogAddress interface{} `db:"BlogAddress" json:"blogAddress"`
	BlogType interface{} `db:"BlogType" json:"blogType"`
	Grade string `db:"Grade" json:"grade"`
	UpdateTime string `db:"UpdateTime" json:"updateTime"`
	Flag int64 `db:"flag" json:"flag"`
}

type Users struct {}

func (this Users) GetAll() *[]User {
	rows := make([]User, 0)
	sql := fmt.Sprintf("select %s from %s", this.columns(), this.table())
	common.DB().Select(&rows, sql)
	return &rows
}

// 获取表名
func (this Users) table() string {
 	return "T_user"
}

//获取数据表的列

func (this Users) columns() string {
	return `Id, Name, QQ, BlogAddress, BlogType, Grade, UpdateTime, flag`
}
