package common

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"database/sql"
	"log"
	"time"
	"errors"
	"manager/models"
)

type NewDb struct {
	account      string
	ipList       []string
	dbName       string
	port         int
	maxIdleConns int
	maxOpenConns int

	db        *sqlx.DB // 连接数据库的db连接
	currentIp string   // 连接数据库的db连接
	lock      sync.Mutex
}

var (
	master *NewDb
)

func DB() *NewDb {
	return master
}

func InitDB() {
	conf := models.Conf()
	master = ConnectDb(conf.DB.Account, conf.DB.Ips, conf.DB.Name,
		conf.DB.Port, conf.DB.MaxIdleConns, conf.DB.MaxOpenConns)
}

func CloseDB() {
	master.Close()
}

func ConnectDb(account string, ipList []string, dbName string, port, maxIdleConns, maxOpenConns int) *NewDb {
	return &NewDb{
		account:           account,
		ipList:            ipList,
		dbName:            dbName,
		port:              port,
		maxIdleConns:      maxIdleConns,
		maxOpenConns:      maxOpenConns,
	}
}

/*
 * 写入一条记录, 默认返回id(多条记录可忽略)
 */
func (this *NewDb) Insert(query string, args ...interface{}) int64 {
	result := this.exec(query, args...)
	id, err := result.LastInsertId()
	this.handleError(err, query, args)
	return id
}

/*
 * 更新记录, 默认返回影响的行数
 */
func (this *NewDb) Update(query string, args ...interface{}) int64 {
	result := this.exec(query, args...)
	rowsCnt, err := result.RowsAffected()
	this.handleError(err, query, args)
	return rowsCnt
}

/*
 * 删除记录, 默认返回影响的行数
 */
func (this *NewDb) Delete(query string, args ...interface{}) int64 {
	result := this.exec(query, args...)
	rowsCnt, err := result.RowsAffected()
	this.handleError(err, query, args)
	return rowsCnt
}

/*
 * 获取一条记录
 */
func (this *NewDb) Get(dest interface{}, query string, args ...interface{}) {
	db := this.getDb(this.db, false)
	err := db.Get(dest, query, args...)
	this.handleError(err, query, args)
}

/*
 * 获取多条记录
 */
func (this *NewDb) Select(dest interface{}, query string, args ...interface{}) {
	db := this.getDb(this.db, false)
	err := db.Select(dest, query, args...)
	this.handleError(err, query, args)
}

/*
 * 执行函数
 */
func (this *NewDb) exec(query string, args ...interface{}) sql.Result {
	if err := recover(); err != nil {
		msg := fmt.Sprintf("[sql: %s][args: %v#][err: %s]", query, args, err)
		log.Printf(msg)
		// TODO 发送邮件
		tips := fmt.Sprintf("数据库错误[err: %s]", err)
		panic(tips)
	}
	db := this.getDb(this.db, false)
	tx := db.MustBegin()
	ret := tx.MustExec(query, args...)
	tx.Commit()

	return ret
}

/*
 * 关闭连接
 */
func (this *NewDb) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("close recovery -> %s\n%s\n", err)
			// TODO 发送堆栈信息邮件至Admin
		}
	}()

	if this.db != nil {
		// 对下面的代码加锁
		this.lock.Lock()
		defer this.lock.Unlock()

		this.db.Close()
		this.db = nil
	}
}

/*
 * 获取数据库
 */
func (this *NewDb) getDb(currentDb *sqlx.DB, retry bool) *sqlx.DB {
	if len(this.ipList) == 0 {
		msg := fmt.Sprintf("数据库[配置]错误: ip列表为空。请检查")
		panic(msg)
	}

	if !retry && this.db != nil {
		return this.db
	}

	// 对下面的代码加锁
	this.lock.Lock()
	defer this.lock.Unlock()

	// 检查是否其他进程是否已经重试连接过数据库
	if this.db != currentDb {
		return this.db
	} else {
		// 释放并重置db变量
		this.Close()

		// 循环3次
		for counter := 0; counter < 3; counter++ {
			for _, ip := range this.ipList {
				// 找到最近连接的一个ip，并从下一个ip开始尝试连接数据库
				if counter == 0 && len(this.currentIp) > 0 {
					// 循环到最近的ip时，置空，避免再次循环
					if this.currentIp == ip {
						this.currentIp = ""
					}
					continue
				}

				// 尝试连接数据库
				db, err := this.open(ip)
				if err == nil {
					this.db = db
					return this.db
				}
			}
			time.Sleep(time.Duration(2^(counter)) * time.Second)
		}
	}
	return this.db
}

func (this *NewDb) open(ip string) (*sqlx.DB, error) {
	openInfo := fmt.Sprintf("%s@tcp(%s:%d)/%s?charset=utf8&timeout=3s", this.account, ip, this.port, this.dbName)
	db, err := sqlx.Connect("mysql", openInfo)
	if err != nil {
		msg := fmt.Sprintf("数据库[%s][连接]错误:%s", openInfo, err.Error())
		log.Println(msg)
		return nil, errors.New(msg)
	}
	db.SetMaxIdleConns(this.maxIdleConns)
	db.SetMaxOpenConns(this.maxOpenConns)
	return db, nil
}

func (this *NewDb) handleError(err error, query string, args ...interface{}) bool {
	isSucc := true
	if err != nil && err != sql.ErrNoRows {
		isSucc = false
		msg := fmt.Sprintf("[sql: %s][args: %v#][err: %s]", query, args, err)
		log.Printf(msg)
		// 发送邮件

		tips := fmt.Sprintf("数据库操作异常[err: %s]", err)
		panic(tips)
	}
	return isSucc
}

