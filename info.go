package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Info struct {
	id   int64
	name string
}

func GetBitcoinsInfo(db *sql.DB) (infos map[string]int64) {
	infos = make(map[string]int64)
	// Execute the query
	rows, err := db.Query("SELECT * FROM cc_info")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		/*
			var value string
			for i, col := range values {
				// Here we can check if the value is nil (NULL value)
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}
				//fmt.Println(columns[i], ": ", value)
			}
		*/
		id, err := strconv.ParseInt(string(values[0]), 10, 64)
		if err != nil {
			fmt.Println("parse info id error:", err)
		} else {
			infos[string(values[2])] = id
		}
		//fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return
}

func AddInfo(name, symbol string, db *sql.DB) (int64, error) {
	insert, err := db.Exec("insert into cc_info (name, symbol, time) values (?, ?, NOW())", name, symbol)
	if err == nil {
		lastid, err := insert.LastInsertId()
		if err == nil {
			return lastid, nil
		}
	}
	return -1, err
}
