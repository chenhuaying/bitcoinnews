package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func AddPrice(id int64, price float64, db *sql.DB) (int64, error) {
	insert, err := db.Exec("insert into cc_price (ccid, price, time) values (?, ?, NOW())", id, price)
	if err == nil {
		lastid, err := insert.LastInsertId()
		if err == nil {
			return lastid, nil
		}
	}
	return -1, err
}
