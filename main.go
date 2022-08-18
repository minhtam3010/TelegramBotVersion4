package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConvertMysqlTimeUnixTime(mysqlTime string) string {
	if mysqlTime == "0" {
		return "0"
	}
	res1 := strings.Replace(mysqlTime, "T", " ", 1)
	res2 := res1[11:19]
	layout := "15:04:05"
	t, err := time.Parse(layout, res2)
	if err != nil {
		panic(err)
	}

	return t.String()[11:19]
}

func ConvertToday(today time.Time) string {
	return today.Format("2006-01-02")
}

type Table1 struct {
	ID          int
	FullName    string
	Active      int
	DateCreated string
	DateUpdated string
}

func CallBotInserted(db *gorm.DB, tableName interface{}) string{
	db.Where("date_created >= DATE(NOW())").Find(&tableName)
	resInsert := "Inserted:\n%v\n-------------------------------------------------------------------------------\n"
	resInsert = fmt.Sprintf(resInsert, tableName)
	return resInsert
}

func CallBotUpdated(db *gorm.DB, tableName interface{}) string {
	db.Where("date_updated >= DATE(NOW()) and active = 1").Find(&tableName)
	resUpdate := "Updated:\n%v\n-------------------------------------------------------------------------------\n"
	resUpdate = fmt.Sprintf(resUpdate, tableName)
	return resUpdate
}

func CallBotDeleted(db *gorm.DB, tableName interface{}) string {
	db.Where("date_updated >= DATE(NOW()) and active = 0").Find(&tableName)
	resDelete := "Deleted:\n%v\n-------------------------------------------------------------------------------\n"

	resDelete = fmt.Sprintf(resDelete, tableName)
	return resDelete
}

func CallTelegramBot(DNS string, BotAPI string, tableName interface{}) {
	db, err := gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		return
	}
	bot, err := tgbotapi.NewBotAPI(BotAPI)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil { // If we got a message

			// splitInputUser := strings.Split(update.Message.Text, " ")
			if strings.ToLower(update.Message.Text) == "info" {

				resInsert := CallBotInserted(db, tableName)
				
				today := "Information about Invoice table: " + ConvertToday(time.Now())
				msgToday := tgbotapi.NewMessage(update.Message.Chat.ID, today)
				
				msgInsert := tgbotapi.NewMessage(update.Message.Chat.ID, resInsert)
				// msgStatisticInsert := tgbotapi.NewMessage(update.Message.Chat.ID, "Total inserted: "+strconv.Itoa(countInvoiceInsert))
				bot.Send(msgToday)
				bot.Send(msgInsert)
				// bot.Send(msgStatisticInsert)

				resUpdate := CallBotUpdated(db, tableName)
				msgUpdate := tgbotapi.NewMessage(update.Message.Chat.ID, resUpdate)
				// msgStatisticUpdate := tgbotapi.NewMessage(update.Message.Chat.ID, "Total updated: "+strconv.Itoa(countInvoiceUpdate))
				bot.Send(msgUpdate)
				// bot.Send(msgStatisticUpdate)

				resDelete := CallBotDeleted(db, tableName)
				msgDelete := tgbotapi.NewMessage(update.Message.Chat.ID, resDelete)

				// msgStatisticDelete := tgbotapi.NewMessage(update.Message.Chat.ID, "Total deleted: "+strconv.Itoa(countInvoiceDelete))

				bot.Send(msgDelete)
				// bot.Send(msgStatisticDelete)
				break
			}
		}
	}
}

func main() {
	dns := "root:leminhtamâđđâ@tcp(127.0.0.1:3306)/BotTest?charset=utf8mb4&parseTime=True&loc=Local"

	var allTable []Table1
	CallTelegramBot(dns, "5423441007:AAGvULEN7X2nZU2uZfoMGUfEQBCrNWzQg7k", allTable)
}
