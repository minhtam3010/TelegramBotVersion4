package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"minhtam/database"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Invoice struct {
	ID                   string
	BalanceID            string
	OrderID              string
	Status               bool
	BalanceBeforeDeposit int
	BalanceAfterDeposit  int
	Note                 string
	Payment              int64
	ClerkID              int
	DeletedUserID        int
	TransactionID        string
	Active               int
	DateCompleted        string
	DateCreated          string
	DateUpdated          string
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

func ConvertToday(today time.Time) string {
	return today.Format("2006-01-02")
}

// func CallTelegramBot(entity struct{})

func CallTelegramBot(BotAPI string, ctx context.Context, tx *sql.Tx) {
	bot, err := tgbotapi.NewBotAPI(BotAPI)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	// log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message

			// splitInputUser := strings.Split(update.Message.Text, " ")
			if strings.ToLower(update.Message.Text) == "info" {

				insertedInvoice, err := tx.QueryContext(ctx, `select id, balance_id, order_id, status, balance_before_deposit, balance_after_deposit, note, payment, 
				clerk_id, deleted_user_id, transaction_id, active, date_completed, date_created, date_updated from invoices where date_created >= date(now())`)
				if err != nil {
					return
				}

				var (
					dataInsert                         Invoice
					allInvoiceInsert                   []interface{}
					dataUpdate                         Invoice
					allInvoiceUpdate                   []interface{}
					allInvoiceDelete                   []interface{}
					dataDelete                         Invoice
					date_completed, date_updated, note sql.NullString
					deleted_user_id                    sql.NullInt32
					status                             int
				)
				for insertedInvoice.Next() {
					err = insertedInvoice.Scan(&dataInsert.ID, &dataInsert.BalanceID, &dataInsert.OrderID, &status, &dataInsert.BalanceBeforeDeposit, &dataInsert.BalanceAfterDeposit, &note, &dataInsert.Payment, &dataInsert.ClerkID, &deleted_user_id, &dataInsert.TransactionID, &dataInsert.Active, &date_completed, &dataInsert.DateCreated, &date_updated)
					if err != nil {
						log.Println(err)
						return
					}

					dataInsert.DateCreated = ConvertMysqlTimeUnixTime(dataInsert.DateCreated)
					if deleted_user_id.Valid {
						dataInsert.DeletedUserID = int(deleted_user_id.Int32)
					}

					if note.Valid {
						dataInsert.Note = note.String
					}
					if date_completed.Valid {
						dataInsert.DateCompleted = ConvertMysqlTimeUnixTime(dataInsert.DateCompleted)
					}
					if date_updated.Valid {
						dataInsert.DateUpdated = ConvertMysqlTimeUnixTime(date_updated.String)
					}
					allInvoiceInsert = append(allInvoiceInsert, dataInsert)
				}
				updatedInvoice, err := tx.QueryContext(ctx, `select id, balance_id, order_id, status, balance_before_deposit, balance_after_deposit, note, 
				payment, clerk_id, deleted_user_id, transaction_id, active, date_completed, date_created, date_updated from invoices where date_updated >= date(now()) and active = 1`)
				if err != nil {
					return
				}

				for updatedInvoice.Next() {

					err = updatedInvoice.Scan(&dataUpdate.ID, &dataUpdate.BalanceID, &dataUpdate.OrderID, &status, &dataUpdate.BalanceBeforeDeposit, &dataUpdate.BalanceAfterDeposit, &note, &dataUpdate.Payment, &dataUpdate.ClerkID, &deleted_user_id, &dataUpdate.TransactionID, &dataUpdate.Active, &date_completed, &dataUpdate.DateCreated, &dataUpdate.DateUpdated)
					if err != nil {
						return
					}

					dataUpdate.DateCreated = ConvertMysqlTimeUnixTime(dataUpdate.DateCreated)
					dataUpdate.DateUpdated = ConvertMysqlTimeUnixTime(dataUpdate.DateUpdated)
					if deleted_user_id.Valid {
						dataUpdate.DeletedUserID = int(deleted_user_id.Int32)
					}

					if note.Valid {
						dataUpdate.Note = note.String
					}
					if date_completed.Valid {
						dataUpdate.DateCompleted = ConvertMysqlTimeUnixTime(dataUpdate.DateCompleted)
					}
					allInvoiceUpdate = append(allInvoiceUpdate, dataUpdate)
				}

				deletedInvoice, err := tx.QueryContext(ctx, `select id, balance_id, order_id, status, balance_before_deposit, balance_after_deposit, note, 
											payment, clerk_id, deleted_user_id, transaction_id, active, date_completed, date_created, date_updated from invoices where date_updated >= date(now()) and active = 0`)
				if err != nil {
					return
				}
				for deletedInvoice.Next() {

					err = deletedInvoice.Scan(&dataDelete.ID, &dataDelete.BalanceID, &dataDelete.OrderID, &status, &dataDelete.BalanceBeforeDeposit, &dataDelete.BalanceAfterDeposit, &note, &dataDelete.Payment, &dataDelete.ClerkID, &deleted_user_id, &dataDelete.TransactionID, &dataDelete.Active, &date_completed, &dataDelete.DateCreated, &dataDelete.DateUpdated)
					if err != nil {
						return
					}

					dataDelete.DateCreated = ConvertMysqlTimeUnixTime(dataDelete.DateCreated)
					dataDelete.DateUpdated = ConvertMysqlTimeUnixTime(dataDelete.DateUpdated)
					if deleted_user_id.Valid {
						dataDelete.DeletedUserID = int(deleted_user_id.Int32)
					}

					if note.Valid {
						dataDelete.Note = note.String
					}
					if date_completed.Valid {
						dataDelete.DateCompleted = ConvertMysqlTimeUnixTime(dataDelete.DateCompleted)
					}
					allInvoiceDelete = append(allInvoiceDelete, dataDelete)
				}

				resInsert := "Inserted:\n"
				countInvoiceInsert := len(allInvoiceInsert)

				for i, data := range allInvoiceInsert {
					invoiceDetail, _ := data.(Invoice)
					n := "Time Created -- " + invoiceDetail.DateCreated
					if i+1 == countInvoiceInsert {
						resInsert += n + ": %v"
					} else {
						resInsert += n + ": %v\n-------------------------------------------------------------------------------\n"
					}
				}
				resUpdate := "Updated:\n"
				countInvoiceUpdate := len(allInvoiceUpdate)

				for i, data := range allInvoiceUpdate {
					invoiceDetail, _ := data.(Invoice)
					n := "Time Updated -- " + invoiceDetail.DateUpdated
					if i+1 == countInvoiceUpdate {
						resUpdate += n + ": %v"
					} else {
						resUpdate += n + ": %v\n-------------------------------------------------------------------------------\n"
					}
				}

				resDelete := "Deleted:\n"
				countInvoiceDelete := len(allInvoiceDelete)

				for i, data := range allInvoiceDelete {
					invoiceDetail, _ := data.(Invoice)
					n := "Time Deleted -- " + invoiceDetail.DateUpdated
					if i+1 == countInvoiceUpdate {
						resDelete += n + ": %v"
					} else {
						resDelete += n + ": %v\n-------------------------------------------------------------------------------\n"
					}
				}
				today := "Information about Invoice table: " + ConvertToday(time.Now())
				msgToday := tgbotapi.NewMessage(update.Message.Chat.ID, today)
				
				resInsert = fmt.Sprintf(resInsert, allInvoiceInsert...)
				msgInsert := tgbotapi.NewMessage(update.Message.Chat.ID, resInsert)
				msgStatisticInsert := tgbotapi.NewMessage(update.Message.Chat.ID, "Total inserted: "+strconv.Itoa(countInvoiceInsert))

				resUpdate = fmt.Sprintf(resUpdate, allInvoiceUpdate...)
				msgUpdate := tgbotapi.NewMessage(update.Message.Chat.ID, resUpdate)
				msgStatisticUpdate := tgbotapi.NewMessage(update.Message.Chat.ID, "Total updated: "+strconv.Itoa(countInvoiceUpdate))
				
				resDelete = fmt.Sprintf(resDelete, allInvoiceDelete...)
				msgDelete := tgbotapi.NewMessage(update.Message.Chat.ID, resDelete)
				msgStatisticDelete := tgbotapi.NewMessage(update.Message.Chat.ID, "Total deleted: "+strconv.Itoa(countInvoiceDelete))

				bot.Send(msgToday)
				bot.Send(msgInsert)
				bot.Send(msgStatisticInsert)
				bot.Send(msgUpdate)
				bot.Send(msgStatisticUpdate)
				bot.Send(msgDelete)
				bot.Send(msgStatisticDelete)
				break
			}

			// switch splitInputUser[0] {
			// case "Insert":
			// 	// ctx := context.WithValue(context.Background(), "tx", tx)
			// 	inserted, err := GetInserted(ctx, tx, splitInputUser[1])
			// 	if err != nil {
			// 		return
			// 	}
			// 	for inserted.Next() {
			// 		err = inserted.Scan(&data.id, &data.lname, &data.fname, &date_created, &date_updated)
			// 		if err != nil {
			// 			return
			// 		}
			// 		data.date_created = ConvertMysqlTimeUnixTime(date_created)
			// 		if date_updated.Valid {
			// 			data.date_updated = ConvertMysqlTimeUnixTime(date_updated.String)
			// 		}
			// 		allTest = append(allTest, data)
			// 	}
			// case "Update":
			// 	// ctx := context.WithValue(context.Background(), "tx", tx)
			// 	updated, err := GetUpdated(ctx, tx, splitInputUser[1])
			// 	if err != nil {
			// 		log.Println(err)

			// 		return
			// 	}
			// 	for updated.Next() {
			// 		err = updated.Scan(&data.id, &data.lname, &data.fname, &date_created, &date_updated)
			// 		if err != nil {
			// 			return
			// 		}
			// 		data.date_created = ConvertMysqlTimeUnixTime(date_created)
			// 		if date_updated.Valid {
			// 			data.date_updated = ConvertMysqlTimeUnixTime(date_updated.String)
			// 		}
			// 		allTest = append(allTest, data)
			// 	}
			// }
			// res := fmt.Sprintf("Inserted: %v", )
			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
			// bot.Send(msg)
			// break
		}
	}
}

func main() {
	tx := database.GetTX()
	ctx := context.WithValue(context.Background(), "tx", tx)
	CallTelegramBot("5423441007:AAGvULEN7X2nZU2uZfoMGUfEQBCrNWzQg7k", ctx, tx)
	// time.Sleep(10 * time.Second)
}
