package main

// origin --> TamTelegram
import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
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

func filterString(myStr string) (res []interface{}, dateTime []string) {
	replacer := strings.NewReplacer("[", "", "{", "", "}", "", "]", "")
	filter := replacer.Replace(myStr)
	split_data := strings.Split(filter, " ")

	for _, str := range split_data {
		re := regexp.MustCompile("2022-[0-9][0-9]-[0-9][0-9]")
		if len(str) == 25 && len(re.FindStringSubmatch(str)) >= 1 {
			re = regexp.MustCompile("[0-9][0-9]:[0-9][0-9]:[0-9][0-9]")
			dateTime = append(dateTime, re.FindStringSubmatch(str)[0])
			res = append(res, re.FindStringSubmatch(str)[0])
		} else {
			res = append(res, str)
		}

	}
	return res, dateTime
}

func CallBotInserted(db *gorm.DB, myStruct []interface{}, tableName []string) []string {
	if len(myStruct) != len(tableName) {
		return []string{errors.New("your input is not equal either the struct or table name").Error()}
	}
	switch len(myStruct) {
	case 0:
		return []string{errors.New("not found any struct or table name in your input").Error()}
	case 1:
		db.Table(tableName[0]).Select("*").Where("date_created >= DATE(NOW())").Find(&myStruct[0])
		resInsert := fmt.Sprintf("%v", myStruct)
		return []string{resInsert}
	default:
		var (
			resInsert []string
		)
		for i, inter := range myStruct {
			db.Table(tableName[i]).Select("*").Where("date_created >= DATE(NOW())").Find(&inter)
			resInsert = append(resInsert, fmt.Sprintf("%v", inter))
		}
		return resInsert
	}
}

func CallBotUpdated(db *gorm.DB, myStruct []interface{}, tableName []string) []string {
	if len(myStruct) != len(tableName) {
		return []string{errors.New("your input is not equal either the struct or table name").Error()}
	}
	switch len(myStruct) {
	case 0:
		return []string{errors.New("not found any struct or table name in your input").Error()}
	case 1:
		db.Table(tableName[0]).Select("*").Where("date_updated >= DATE(NOW()) and active = 1").Find(&myStruct[0])
		resUpdate := fmt.Sprintf("%v", myStruct)
		return []string{resUpdate}
	default:
		var (
			resUpdate []string
		)
		for i, inter := range myStruct {
			db.Table(tableName[i]).Select("*").Where("date_updated >= DATE(NOW()) and active = 1").Find(&inter)
			resUpdate = append(resUpdate, fmt.Sprintf("%v", inter))
		}
		return resUpdate
	}
}

func CallBotDeleted(db *gorm.DB, myStruct []interface{}, tableName []string) []string {
	if len(myStruct) != len(tableName) {
		return []string{errors.New("your input is not equal either the struct or table name").Error()}
	}
	switch len(myStruct) {
	case 0:
		return []string{errors.New("not found any struct or table name in your input").Error()}
	case 1:
		db.Table(tableName[0]).Select("*").Where("date_updated >= DATE(NOW()) and active = 0").Find(&myStruct[0])
		resDelete := fmt.Sprintf("%v", myStruct)
		return []string{resDelete}
	default:
		var (
			resDelete []string
		)
		for i, inter := range myStruct {
			db.Table(tableName[i]).Select("*").Where("date_updated >= DATE(NOW()) and active = 0").Find(&inter)
			resDelete = append(resDelete, fmt.Sprintf("%v", inter))
		}
		return resDelete
	}
}

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

				today := "Automation checking CRUD table: " + ConvertToday(time.Now())
				msgToday := tgbotapi.NewMessage(update.Message.Chat.ID, today)
				bot.Send(msgToday)
				getAllInsertData := CallBotInserted(db, myStruct, tableName)
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
						str, datetime = filterString(i)
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
					getAllUpdateData := CallBotUpdated(db, myStruct, tableName)
					split_data = strings.Split(getAllUpdateData[i], "} {")
					for _, i := range split_data {
						str, datetime = filterString(i)
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
					getAllDeleteData := CallBotDeleted(db, myStruct, tableName)
					split_data = strings.Split(getAllDeleteData[i], "} {")
					for _, i := range split_data {
						str, datetime = filterString(i)
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

func main() {
	dns := "root:leminhtamâđđâ@tcp(127.0.0.1:3306)/BotTest?charset=utf8mb4&parseTime=True&loc=Local"

	// var allTable []interface{}
	var allTable1 []T1
	var allTable2 []T2
	var allTable3 []T3

	allTable := []interface{}{
		allTable1,
		allTable2,
		allTable3,
	}
	allTableName := []string{"table1", "table2", "table3"}

	CallTelegramBot(dns, "5423441007:AAGvULEN7X2nZU2uZfoMGUfEQBCrNWzQg7k", allTable, allTableName)
}
