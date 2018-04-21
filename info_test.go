package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestGetInfo(t *testing.T) {
	// Open database connection
	db, err := sql.Open("mysql", "chy:123456@tcp(192.168.56.102:3306)/test")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	infos := GetBitcoinsInfo(db)
	fmt.Println(infos)
	chyId, ok := infos["chy"]
	if ok {
		t.Error("no data in mysql but find")
	} else {
		fmt.Println("id = ", chyId)
	}

	btcId, ok := infos["BTC"]
	if ok {
		fmt.Println("btcId = ", btcId)
	} else {
		t.Error("data in mysql but not find")
	}
}

func TestAddInfo(t *testing.T) {
	// Open database connection
	db, err := sql.Open("mysql", "chy:123456@tcp(192.168.56.102:3306)/test")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	name, err := AddInfo("Bitcoin", "BTC", db)
	if err != nil {
		t.Error("add info error:", err)
	} else {
		fmt.Println(name)
	}
}
