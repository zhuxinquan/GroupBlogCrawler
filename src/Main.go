package main

import (
	"common"
	"manager/models"
	"manager"
	"time"
)

func main(){
	models.InitConfig()
	common.InitDB()
	for true {
		manager.Manager()
		time.Sleep(3 * time.Hour)
	}
	common.CloseDB()
}
