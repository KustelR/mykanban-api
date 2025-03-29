package db_driver

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetDb(connString string) *sql.DB {
	var err error
	var db *sql.DB

	failCount := 0
	loopControl := true
	for loopControl {
		if failCount > 0 {
			time.Sleep(time.Second * 10)
		}
		db, err = sql.Open("mysql", connString)
		if err != nil {
			failCount += 1
			log.Printf("can't connect to database, err: %s", err)
			continue
		}
		err = db.Ping()
		if err != nil {
			failCount += 1
			log.Printf("failed pinging database: %s", err)
			continue
		}
		loopControl = false
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
