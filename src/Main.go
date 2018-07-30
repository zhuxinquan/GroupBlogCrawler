package main

import (
	"common"
	"manager"
	"manager/models"
	"time"
)

func main() {
	models.InitConfig()
	common.InitDB()
	for true {
		manager.Manager()
		time.Sleep(3 * time.Hour)
	}
	common.CloseDB()
}
