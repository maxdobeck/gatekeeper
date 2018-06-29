// Package models holds all the data layer interfaces and transactions
package models

import (
	"database/sql"
	_ "github.com/lib/pq" // github.com/lib/pq
	"log"
)

var db *sql.DB

// ConnToDB connects the database.  This should be called at the app start
func ConnToDB(dbURL string) {
	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Println(err)
	}
	if err = db.Ping(); err != nil {
		log.Println(err)
	}
}
