package models

import (
	// _ "github.com/lib/pq" // github.com/lib/pq
	"fmt"
	"os"
	"testing"
)

// TestCreateShift creates a new shift
func TestCreateShift(t *testing.T) {
	ConnToDB(os.Getenv("PGURL"))

	_, delErr := Db.Query("DELETE FROM schedules WHERE title like 'Test Shift'")
	fmt.Println(delErr)

	spoofShifts()

	rows, errors := Db.Query("SELECT id FROM members LIMIT 1;")
	if errors != nil {
		fmt.Println(errors)
	}
	defer rows.Close()
	var memberID string
	for rows.Next() {
		err := rows.Scan(&memberID)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Member used for testing: ", memberID)
	}

	srows, serrors := Db.Query("SELECT id FROM shifts LIMIT 1;")
	if serrors != nil {
		fmt.Println(serrors)
	}
	defer srows.Close()
	var scheduleID string
	for srows.Next() {
		rerr := srows.Scan(&scheduleID)
		if rerr != nil {
			fmt.Println(rerr)
		}
		fmt.Println("Schedule ID used for test Shift: ", scheduleID)
	}

	s := Shift{
		Title:        "Test Shift",
		Schedule:     scheduleID,
		Start:        "715",
		End:          "1200",
		MinEnrollees: "1",
		Days:         [7]string{"", "Monday", "Tuesday", "", "", "Friday", ""},
	}

	var newShiftError error
	newShiftError = CreateShift(&s)
	if newShiftError != nil {
		fmt.Println(newShiftError)
		t.Fail()
	}

	var record string
	err := Db.QueryRow("SELECT title FROM schedules WHERE title LIKE 'Test Shift'").Scan(&record)
	if err != nil {
		fmt.Println("Test Failed because: ", err)
		t.Fail()
	}
	if record != "Test Shift" {
		t.Fail()
	}
	cleanupShifts()
}

// Helpers for shift tests
func spoofShifts() {
	m := NewMember{
		Name:      "Test Member",
		Email:     "testuser44@gmail.com",
		Email2:    "testuser44@gmail.com",
		Password:  "superduper",
		Password2: "superduper",
	}
	if CreateMember(&m) != nil {
		fmt.Println("Member may already be there")
	}

	l := make([]*Schedule, 1)
	l[0] = &Schedule{"", "My Morning Shift", GetMemberID("testuser44@gmail.com")}

	for i := range l {
		if CreateSchedule(l[i]) != nil {
			fmt.Println("Schedule may already exist")
		}
	}
}

// cleanupShifts undoes the shift spoofing func with a cascade delete
func cleanupShifts() {
	_, err := Db.Query("DELETE FROM members WHERE email LIKE 'testuser44@gmail.com'")
	if err != nil {
		fmt.Println(err)
	}
}
