package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"minhtam/database"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Test1 struct {
	id           int
	lname        string
	fname        string
	date_created string
	date_updated string
}

func ConvertMysqlTimeUnixTime(mysqlTime string) string {
	if mysqlTime == "0" {
		return "0"
	}
	res1 := strings.Replace(mysqlTime, "T", " ", 1)
	res2 := res1[11:19]
	// YYYY-MM-DD
	layout := "15:04:05"
	t, err := time.Parse(layout, res2)
	if err != nil {
		panic(err)
	}

	return t.String()[11:19]
}

func GetInserted(ctx context.Context, tx *sql.Tx, tableName string) (rows *sql.Rows, err error) {
	rows, err = tx.QueryContext(ctx, "select * from "+tableName+" where date_created >= date(now())")
	if err != nil {
		return
	}

	return rows, err
}

func GetUpdated(ctx context.Context, tx *sql.Tx, tableName string) (rows *sql.Rows, err error) {
	rows, err = tx.QueryContext(ctx, "select * from "+tableName+" where date_updated >= date(now()) and active = 1")
	if err != nil {
		return
	}

	return rows, err
}

func main() {
	bot, err := tgbotapi.NewBotAPI("5423441007:AAGvULEN7X2nZU2uZfoMGUfEQBCrNWzQg7k")
	if err != nil {
		log.Println("Problem")
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var (
		allTest      []Test1
		data         Test1
		date_created string
		date_updated sql.NullString
	)
	for update := range updates {
		if update.Message != nil { // If we got a message

			splitInputUser := strings.Split(update.Message.Text, " ")
			switch splitInputUser[0] {
			case "Insert":
				tx := database.GetTX()

				ctx := context.WithValue(context.Background(), "tx", tx)
				inserted, err := GetInserted(ctx, tx, splitInputUser[1])
				if err != nil {
					return
				}
				for inserted.Next() {
					err = inserted.Scan(&data.id, &data.lname, &data.fname, &date_created, &date_updated)
					if err != nil {
						return
					}
					data.date_created = ConvertMysqlTimeUnixTime(date_created)
					if date_updated.Valid {
						data.date_updated = ConvertMysqlTimeUnixTime(date_updated.String)
					}
					allTest = append(allTest, data)
				}
			case "Update":
				tx := database.GetTX()

				ctx := context.WithValue(context.Background(), "tx", tx)
				updated, err := GetUpdated(ctx, tx, splitInputUser[1])
				if err != nil {
					log.Println(err)

					return
				}
				for updated.Next() {
					err = updated.Scan(&data.id, &data.lname, &data.fname, &date_created, &date_updated)
					if err != nil {
						return
					}
					data.date_created = ConvertMysqlTimeUnixTime(date_created)
					if date_updated.Valid {
						data.date_updated = ConvertMysqlTimeUnixTime(date_updated.String)
					}
					allTest = append(allTest, data)
				}
			}
			res := fmt.Sprintf("Inserted: %v", allTest)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
			bot.Send(msg)
		}
	}
}
