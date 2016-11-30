package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func GetDB() *sql.DB {
	if db != nil {
		return db
	}

	var err error
	db, err = sql.Open("mysql", "root:@/hogwarts")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
