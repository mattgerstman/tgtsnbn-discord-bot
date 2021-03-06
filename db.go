package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var internalDB *sql.DB

/**
 * Returns a pointer to the database.
 * Initializes the database if necessary.
 */
func GetDB() *sql.DB {
	if internalDB != nil {
		return internalDB
	}

	config := GetConfig()

	var err error
	internalDB, err = sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// Sanity check that the internalDB is up. If it's not quit.
	err = internalDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return internalDB
}
