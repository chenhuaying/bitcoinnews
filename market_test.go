package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestAddMarket(t *testing.T) {
	// Open database connection
	db, err := sql.Open("mysql", "chy:123456@tcp(192.168.56.102:3306)/test")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	m, err := AddMarketCap(8, 10000000000, 100001, 201, -0.88, db)
	if err != nil {
		t.Error("add AddMarketCap error:", err)
	} else {
		fmt.Println(m)
	}
}
