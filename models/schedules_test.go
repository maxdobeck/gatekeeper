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

	s := Schedule{
		Title:   "Test Schedule",
		OwnerID: memberId,
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
	cleanupDb()
}

// TestCreateScheduleForNonexistentUser will try and create a schedule for a user (owner) that doesn't exist
/*TestCreateScheduleForNonexistentUser(t *testing.T) {
}*/

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
	cleanupDb()
}

// TestGetSchedules trys to get all of a member's schedules
func TestGetSchedules(t *testing.T) {
	populateDb()
	memberId := GetMemberID("testuser33@gmail.com")
	if memberId == "" {
		t.Fail()
	}

	var s []Schedule
	s, getAllErr := GetSchedules(memberId)
	if getAllErr != nil {
		fmt.Println("All schedules: ", s)
		t.Fail()
	}
	cleanupDb()
}

// TestGetScheduleById attempts to get a schedule by the ID of the schedule
func TestGetScheduleById(t *testing.T) {
	populateDb()
	var s string
	err := Db.QueryRow("SELECT id FROM schedules LIMIT 1").Scan(&s)
	if err != nil {
		fmt.Println("Could not find schedule: ", s)
		t.Fail()
	}
	schedule, err := GetScheduleById(s)
	if err != nil {
		fmt.Println("Could not find schedule: ", s)
		t.Fail()
	}
	if schedule.Id != s {
		fmt.Printf("Could not find schedule. Target schedule id %s != record from DB %s", s, schedule)
		t.Fail()
	}
	cleanupDb()
}

func TestDeleteSchedule(t *testing.T) {
	populateDb()
	var s string
	var err error
	err = Db.QueryRow("SELECT id FROM schedules LIMIT 1").Scan(&s)
	if err != nil {
		fmt.Println("Could not find a schedule to test on")
		t.Fail()
	}
	err = DeleteSchedule(s)
	if err != nil {
		fmt.Println("Could not delete schedule: ", s)
		t.Fail()
	}
	cleanupDb()
}

// Helpers
func populateDb() {
	m := NewMember{
		Name:      "Test Member",
		Email:     "testuser33@gmail.com",
		Email2:    "testuser33@gmail.com",
		Password:  "superduper",
		Password2: "superduper",
	}
	if CreateMember(&m) != nil {
		fmt.Println("Member may already be there")
	}

	l := make([]*Schedule, 4)
	l[0] = &Schedule{"", "Test Test Schedule", GetMemberID("testuser33@gmail.com")}
	l[1] = &Schedule{"", "My 2nd Schedule", GetMemberID("testuser33@gmail.com")}
	l[2] = &Schedule{"", "My 3rd Schedule", GetMemberID("testuser33@gmail.com")}
	l[3] = &Schedule{"", "My 4th Schedule", GetMemberID("testuser33@gmail.com")}

	for i := range l {
		if CreateSchedule(l[i]) != nil {
			fmt.Println("Schedule may already exist")
		}
	}
}

// cleanupDb undoes the populateDb
func cleanupDb() {
	_, err := Db.Query("DELETE FROM members WHERE email LIKE 'testuser33@gmail.com'")
	if err != nil {
		fmt.Println(err)
	}
}
