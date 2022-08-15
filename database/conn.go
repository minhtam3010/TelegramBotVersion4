package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	user     = "root"
	password = "quynhnhu2010"
	host     = "localhost"
	port     = "3306"
	name     = "ieltscenter"
	sslmode  = "false"
)

func OpenConnection() *sql.Tx {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&tls=%s&timeout=30s",
		user, password, host, port, name, sslmode,
	))

	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	return tx
}

func GetTX() *sql.Tx {
	tx := OpenConnection()
	return tx
}
