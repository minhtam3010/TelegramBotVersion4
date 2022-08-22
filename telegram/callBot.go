package telegram

import (
	"fmt"
	"log"
	"minhtam/convert"
	"minhtam/database"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

func CallTelegramBot(DNS string, BotAPI string, myStruct []interface{}, tableName []string) {
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

				today := "Automation checking CRUD table: " + convert.ConvertToday(time.Now())
				msgToday := tgbotapi.NewMessage(update.Message.Chat.ID, today)
				bot.Send(msgToday)
				getAllInsertData := database.GetData(db, myStruct, tableName, "insert")
				for i := 0; i < len(myStruct); i++ {
					var (
						str             []interface{}
						resInsert       []interface{}
						resUpdate       []interface{}
						resDelete       []interface{}
						datetime        []string
						dateTime_insert [][]string
						dateTime_update [][]string
						dateTime_delete [][]string
					)
					split_data := strings.Split(getAllInsertData[i], "} {")
					for _, i := range split_data {
						str, datetime = convert.FilterString(i)
						resInsert = append(resInsert, str)
						dateTime_insert = append(dateTime_insert, datetime)
					}

					line := "---------------------------------------------------------------------------"
					msgLine := tgbotapi.NewMessage(update.Message.Chat.ID, line)
					bot.Send(msgLine)
					fmt.Println(tableName[i])
					info := "Information about the '" + strings.ToUpper(tableName[i]) + "' table"
					msgInfo := tgbotapi.NewMessage(update.Message.Chat.ID, info)
					bot.Send(msgInfo)

					resInserted := "Inserted:\n"
					res := fmt.Sprintf("%v", resInsert[0])
					if res != "[]" {
						for i, _ := range resInsert {
							getDate := dateTime_insert[i][0]
							n := "Time Created -- " + getDate
							if i+1 == len(resInsert) {
								resInserted += n + ": %v"
							} else {
								resInserted += n + ": %v\n-------------------------------------------------------------------------------\n"
							}
						}
						resInserted = fmt.Sprintf(resInserted, resInsert...)
						msgInsert := tgbotapi.NewMessage(update.Message.Chat.ID, resInserted)
						bot.Send(msgInsert)
						msgStatisticInsert := tgbotapi.NewMessage(update.Message.Chat.ID, "Total inserted: "+strconv.Itoa(len(resInsert)))
						bot.Send(msgStatisticInsert)
					} else {
						msgInsert := tgbotapi.NewMessage(update.Message.Chat.ID, "None of data inserted today!!!")
						bot.Send(msgInsert)
					}
					getAllUpdateData := database.GetData(db, myStruct, tableName, "update")
					split_data = strings.Split(getAllUpdateData[i], "} {")
					for _, i := range split_data {
						str, datetime = convert.FilterString(i)
						resUpdate = append(resUpdate, str)
						dateTime_update = append(dateTime_update, datetime)
					}

					resUpdated := "Updated:\n"
					res = fmt.Sprintf("%v", resUpdate[0])
					if res != "[]" {
						for i, _ := range resUpdate {
							getDate := dateTime_update[i][1]
							n := "Time Updated -- " + getDate
							if i+1 == len(resUpdate) {
								resUpdated += n + ": %v"
							} else {
								resUpdated += n + ": %v\n-------------------------------------------------------------------------------\n"
							}
						}
						resUpdated = fmt.Sprintf(resUpdated, resUpdate...)
						msgUpdate := tgbotapi.NewMessage(update.Message.Chat.ID, resUpdated)
						msgStatisticUpdate := tgbotapi.NewMessage(update.Message.Chat.ID, "Total updated: "+strconv.Itoa(len(resUpdate)))
						bot.Send(msgUpdate)
						bot.Send(msgStatisticUpdate)
					} else {
						msgUpdate := tgbotapi.NewMessage(update.Message.Chat.ID, "None of data updated today!!!")
						bot.Send(msgUpdate)
					}
					getAllDeleteData := database.GetData(db, myStruct, tableName, "delete")
					split_data = strings.Split(getAllDeleteData[i], "} {")
					for _, i := range split_data {
						str, datetime = convert.FilterString(i)
						resDelete = append(resDelete, str)
						dateTime_delete = append(dateTime_delete, datetime)
					}

					resDeleted := "Deleted:\n"
					res = fmt.Sprintf("%v", resDelete[0])
					if res != "[]" {
						for i, _ := range resDelete {
							getDate := dateTime_delete[i][1]
							n := "Time Deleted -- " + getDate
							if i+1 == len(resDelete) {
								resDeleted += n + ": %v"
							} else {
								resDeleted += n + ": %v\n-------------------------------------------------------------------------------\n"
							}
						}
						resDeleted = fmt.Sprintf(resDeleted, resDelete...)
						msgDelete := tgbotapi.NewMessage(update.Message.Chat.ID, resDeleted)

						msgStatisticDelete := tgbotapi.NewMessage(update.Message.Chat.ID, "Total deleted: "+strconv.Itoa(len(resDelete)))

						bot.Send(msgDelete)
						bot.Send(msgStatisticDelete)
					} else {
						msgDelete := tgbotapi.NewMessage(update.Message.Chat.ID, "None of data deleted today!!!")
						bot.Send(msgDelete)
					}
				}
			}
			break
		}
	}
}
