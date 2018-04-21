package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func AddMarketCap(id, marketcap, circulating_supply, volume_24h int64, change_24h float64, db *sql.DB) (int64, error) {
	insert, err := db.Exec("insert into cc_market (ccid, marketcap, circulating_supply, volume_24h, change_24h, time) values (?, ?, ?, ?, ?, NOW())",
		id, marketcap, circulating_supply, volume_24h, change_24h)
	if err == nil {
		lastid, err := insert.LastInsertId()
		if err == nil {
			return lastid, nil
		}
	}
	return -1, err
}
