package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"minhtam/singleton"
	"minhtam/telegram"
	"time"

	driver "github.com/go-sql-driver/mysql"
)

// origin --> TamTelegram

type ClassUser struct {
	UserClassID   string `json:"user_class_id"`
	UserID        int    `json:"user_id"`
	ClassID       string `json:"class_id"`
	IsWaitingRoom bool   `json:"is_waiting"`
	Active        int    `json:"active"`
	CreatedUserID int    `json:"created_user_id"`
	DeletedUserID int    `json:"delete_user_id"`
	DateCreated   string `json:"date_created"`
	DateUpdated   string `json:"date_updated"`
}

type Invoice struct {
	ID                   string  `json:"id"`
	BalanceID            string  `json:"balance_id"`
	Title                string  `json:"title"`
	Status               bool    `json:"status"`
	BalanceBeforeDeposit float32 `json:"balance_before_deposit"`
	BalanceAfterDeposit  float32 `json:"balance_after_deposit"`
	Note                 string  `json:"note"`
	Payment              float32 `json:"payment"`
	ClerkID              int     `json:"clerk_id"`
	DeletedUserID        int     `json:"deleted_user_id"`
	TransactionID        string  `json:"transaction_id"`

	Active      int    `json:"active"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

type Order struct {
	ID                     string  `json:"id"`
	Status                 bool    `json:"status"`
	Amount                 float32 `json:"amount"`
	Discount               float32 `json:"discount"`
	AssignDiscountPersonID int     `json:"person_id_assigning_discount"`
	Note                   string  `json:"note"`
	BackNote               string  `json:"back_note"`
	CreatedUserID          int     `json:"created_user_id"`
	Active                 int     `json:"active"`
	DateCreated            string  `json:"date_created"`
	DateUpdated            string  `json:"date_updated"`
}

func main() {
	singleton.InitConfig(".")

	if singleton.Cfg.SSLMode != "false" {
		certBytes, err := base64.StdEncoding.DecodeString(singleton.Cfg.CACERTBASE64)
		if err != nil {
			log.Fatalf("unable to read in the cert file: %s", err)
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(certBytes); !ok {
			log.Fatal("failed to parse sql CA-Cert")
		}

		tlsConfig := &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}

		if err := driver.RegisterTLSConfig("custom", tlsConfig); err != nil {
			panic(err)
		}
	}

	gormCustomMySQLCfg := driver.Config{
		User:   singleton.Cfg.User,
		Passwd: singleton.Cfg.Password,
		Addr:   fmt.Sprintf("%s:%s", singleton.Cfg.DatabaseHost, singleton.Cfg.DatabasePort), //IP:PORT
		Net:    "tcp",
		DBName: singleton.Cfg.Name,
		Loc:    time.Local,
		// AllowNativePasswords: true,
		// Params:               o,
	}

	gormCustomMySQLCfg.TLSConfig = "custom"
	str := gormCustomMySQLCfg.FormatDSN()

	var (
		user_class []ClassUser
		invoices   []Invoice
		orders     []Order
	)

	allTable := []interface{}{
		user_class,
		invoices,
		orders,
	}
	allTableName := []string{"class_user", "invoices", "orders"}

	telegram.CallTelegramBot(str, "5696295625:AAGx8rx_NFmZuv8zflNJqu17y8l8aSilB7A", allTable, allTableName)

	return
}
