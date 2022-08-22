package main

import (
	"minhtam/telegram"
)

// origin --> TamTelegram

type T1 struct {
	ID          int
	FullName    string
	Active      int
	DateCreated string
	DateUpdated string
}

type T2 struct {
	ID          int
	Description string
	Parent      string
	BOD         string
	Active      int
	DateCreated string
	DateUpdated string
}

type T3 struct {
	ID          int
	Parent      string
	Active      int
	DateCreated string
	DateUpdated string
}

func main() {
	dns := "root:quynhnhu2010@tcp(127.0.0.1:3306)/ieltscenter?charset=utf8mb4&parseTime=True&loc=Local"

	var allTable1 []T1
	var allTable2 []T2
	var allTable3 []T3

	allTable := []interface{}{
		allTable1,
		allTable2,
		allTable3,
	}
	allTableName := []string{"table1", "table2", "table3"}

	telegram.CallTelegramBot(dns, "5705043100:AAHfOypgmWT0ICop_VgbLfHmT6QvpLCj1-4", allTable, allTableName)
}
