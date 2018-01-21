package gatekeeper

import (
	"fmt"
	"database/sql"
	"os"
)

func getMemberID (email string) (memberID string) {
	connStr := os.Getenv("PGURL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}
	sqlErr := db.QueryRow("SELECT id FROM members WHERE email = $1", email).Scan(&memberID)
	if sqlErr == sql.ErrNoRows {
		memberID = ""
		return
	}
	if sqlErr != nil {
		fmt.Println(sqlErr)
	}
	return
}
