package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() *sql.DB {
	if db != nil {
		return db
	}

	var err error
	db, err = sql.Open("mysql", "root:@/hogwarts")
	if err != nil {
		log.Fatal(err)
	}

	// Sanity check that the db is up. If it's not quit.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
