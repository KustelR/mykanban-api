package db_driver

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetDb(connString string) *sql.DB {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(fmt.Errorf("can't connect to database, err: %s", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("failed pinging database: %s", err))
	}
	fmt.Println("Connected to MySql DB")
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetConnMaxIdleTime(time.Minute * 3)
	return db
}

type NotFoundError struct {
	thing string
	query *string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s was not found", e.thing)
}
