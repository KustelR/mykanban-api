package db_driver

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
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

func GetMaxDrawOrder(cols []string, data []sql.RawBytes) (int, error) {
	for idx, item := range data {
		if len(item) <= 0 {
			return 0, nil
		}
		if cols[idx] == "max(draw_order)" {
			lastOrder, err := strconv.Atoi(string(item))
			if err != nil {
				return -1, err
			}
			return lastOrder, nil
		}
	}
	return -1, fmt.Errorf("provided query does not contained draw_order")
}
