package telegram

import (
	"fmt"
	"log"
	"minhtam/PDF"
	"minhtam/convert"
	"minhtam/database"
	"os"
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

			if strings.ToLower(update.Message.Text) == "info" {

				today := "Automation checking CRUD table: " + convert.ConvertToday(time.Now())
				msgToday := tgbotapi.NewMessage(update.Message.Chat.ID, today)
				bot.Send(msgToday)
				getAllInsertData := database.GetData(db, myStruct, tableName, "insert")
				var (
					total   [][]string
					crud    = []string{"Inserted", "Updated", "Deleted"}
					res_all []string
				)
				for i := 0; i < len(myStruct); i++ {
					var (
						total_crud      []string
						str             []interface{}
						resInsert       []interface{}
						resUpdate       []interface{}
						resDelete       []interface{}
						datetime        []string
						dateTime_insert [][]string
						dateTime_update [][]string
						dateTime_delete [][]string
						sum_crud        int
					)
					total_crud = append(total_crud, strings.ToUpper(tableName[i]))
					split_data := strings.Split(getAllInsertData[i], "} {")
					for _, i := range split_data {
						str, datetime = convert.FilterString(i)
						resInsert = append(resInsert, str)
						dateTime_insert = append(dateTime_insert, datetime)
					}

					resInserted := ""
					res := fmt.Sprintf("%v", resInsert[0])
					if res != "[]" {
						for i, _ := range resInsert {
							getDate := dateTime_insert[i][0]
							n := "Time Created -- " + getDate
							if i+1 == len(resInsert) {
								resInserted += "+ " + n + ": %v"
							} else {
								resInserted += "+ " + n + ": %v\n\n"
							}
						}
						total_crud = append(total_crud, strconv.Itoa(len(resInsert)))
						resInserted = fmt.Sprintf(resInserted, resInsert...)
						res_all = append(res_all, resInserted)
						sum_crud += len(resInsert)
					} else {
						total_crud = append(total_crud, "0")
						res_all = append(res_all, "")
						sum_crud += 0
					}
					getAllUpdateData := database.GetData(db, myStruct, tableName, "update")
					split_data = strings.Split(getAllUpdateData[i], "} {")
					for _, i := range split_data {
						str, datetime = convert.FilterString(i)
						resUpdate = append(resUpdate, str)
						dateTime_update = append(dateTime_update, datetime)
					}

					resUpdated := ""
					res = fmt.Sprintf("%v", resUpdate[0])
					if res != "[]" {
						for i, _ := range resUpdate {
							getDate := dateTime_update[i][1]
							n := "Time Updated -- " + getDate
							if i+1 == len(resUpdate) {
								resUpdated += "+ " + n + ": %v"
							} else {
								resUpdated += "+ " + n + ": %v\n\n"
							}
						}
						total_crud = append(total_crud, strconv.Itoa(len(resUpdate)))
						resUpdated = fmt.Sprintf(resUpdated, resUpdate...)
						res_all = append(res_all, resUpdated)
						sum_crud += len(resUpdate)
					} else {
						res_all = append(res_all, "")
						total_crud = append(total_crud, "0")
						sum_crud += 0
					}
					getAllDeleteData := database.GetData(db, myStruct, tableName, "delete")
					split_data = strings.Split(getAllDeleteData[i], "} {")
					for _, i := range split_data {
						str, datetime = convert.FilterString(i)
						resDelete = append(resDelete, str)
						dateTime_delete = append(dateTime_delete, datetime)
					}

					resDeleted := ""
					res = fmt.Sprintf("%v", resDelete[0])
					if res != "[]" {
						for i, _ := range resDelete {
							getDate := dateTime_delete[i][1]
							n := "Time Deleted -- " + getDate
							if i+1 == len(resDelete) {
								resDeleted += "+ " + n + ": %v"
							} else {
								resDeleted += "+ " + n + ": %v\n\n"
							}
						}
						total_crud = append(total_crud, strconv.Itoa(len(resDelete)))
						resDeleted = fmt.Sprintf(resDeleted, resDelete...)
						res_all = append(res_all, resDeleted)
						sum_crud += len(resDelete)
					} else {
						total_crud = append(total_crud, "0")
						res_all = append(res_all, "")
						sum_crud += 0
					}

					total_crud = append(total_crud, strconv.Itoa(sum_crud))
					// if convert.Sum(total) != 0 {
					// 	dashboard.BuildChart(m, total, tableName[i])
					// } else {
					// 	fmt.Println("Nothing happened today")
					// }
					total = append(total, total_crud)
				}
				err = PDF.CreatePDF(total, tableName, crud, res_all)
				if err != nil {
					panic(err)
				}
				f, err := os.ReadFile("report.pdf")
				if err != nil {
					panic(err)
				}

				FileBytes := tgbotapi.FileBytes{
					Name:  "report.pdf",
					Bytes: f,
				}
				msg := tgbotapi.NewDocument(update.Message.Chat.ID, FileBytes)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				fmt.Println("PDF saved successfully")
			}
			break
		}
	}
}
