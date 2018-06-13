package manager

import (
	"fmt"
	"manager/models/dbs"
)

// 负责抓取的调用
func Manager() {
	users := dbs.Users{}.GetAll()
	for _, user := range *users {
		if user.Flag == 0 {
			continue
		}
		if user.BlogAddress != nil && user.BlogType != nil {
			switch fmt.Sprintf("%s", user.BlogType) {
			case "csdn":
				Crawlers{}.GetCsdn(fmt.Sprintf("%s", user.BlogAddress), user.Id)
				break
			case "wordpress":
				Crawlers{}.GetWordpress(fmt.Sprintf("%s", user.BlogAddress), user.Id)
				break
			}
		}
	}
}
