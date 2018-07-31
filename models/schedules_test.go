package models

import (
	// _ "github.com/lib/pq" // github.com/lib/pq
	"fmt"
	"os"
	"testing"
)

// TestCreateSchedule will try and create a schedule using an existing user's id
func TestCreateSchedule(t *testing.T) {
	ConnToDB(os.Getenv("PGURL"))

	_, delErr := Db.Query("DELETE FROM schedules WHERE title like 'Test Schedule'")
	fmt.Println(delErr)

	populateDb()

	rows, errors := Db.Query("SELECT id FROM members LIMIT 1;")
	if errors != nil {
		fmt.Println(errors)
	}
	defer rows.Close()
	var memberId string
	for rows.Next() {
		err := rows.Scan(&memberId)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(memberId)
	}

	s := NewSchedule{
		Title: "Test Schedule",
		Owner: memberId,
	}

	var newScheduleError error
	newScheduleError = CreateSchedule(&s)
	if newScheduleError != nil {
		fmt.Println(newScheduleError)
		t.Fail()
	}

	var record string
	err := Db.QueryRow("SELECT title FROM schedules WHERE title LIKE 'Test Schedule'").Scan(&record)
	if err != nil {
		fmt.Println("Test Failed because: ", err)
		t.Fail()
	}
	if record != "Test Schedule" {
		t.Fail()
	}

}

// TestUpdateTitle will change the Title of a schedule
func TestUpdateTitle(t *testing.T) {
	populateDb()
	var scheduleId string
	err := Db.QueryRow("SELECT id FROM schedules LIMIT 1").Scan(&scheduleId)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	updateErr := UpdateScheduleTitle(scheduleId, "New Schedule Title Added")
	if updateErr != nil {
		fmt.Println("Failed to update schedule title: ", updateErr)
		t.Fail()
	}
}

// Helper to populate the DB with reliable data
func populateDb() {
	m := NewMember{
		Name:      "Test Member",
		Email:     "testtest@gmail.com",
		Email2:    "testtest@gmail.com",
		Password:  "superduper",
		Password2: "superduper",
	}

	if CreateMember(&m) != nil {
		fmt.Println("Member may already be there")
	}

	s := NewSchedule{
		Title: "Test Test Schedule",
		Owner: GetMemberID("testtest@gmail.com"),
	}

	if CreateSchedule(&s) != nil {
		fmt.Println("Schedule may already exist")
	}
}
