package models

import (
	// _ "github.com/lib/pq" // github.com/lib/pq
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

// TestCreateShift creates a new shift
func TestCreateShift(t *testing.T) {
	ConnToDB(os.Getenv("PGURL"))

	_, delErr := Db.Query("DELETE FROM schedules WHERE title like 'Test Shift'")
	log.Info(delErr)

	spoofShifts()
	memberID := GetMemberID("testuser144@gmail.com")
	log.Info("Member ID used for testing: ", memberID)

	var scheduleID string
	schedules, getSchedErr := GetSchedules(memberID)
	if getSchedErr != nil {
		t.Error("Error getting schedule.", getSchedErr)
		t.FailNow()
	}
	scheduleID = schedules[0].ID
	log.Info("Schedule ID used for test Shift: ", scheduleID)

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
		log.Info("Test failed while creating new shift: ", newShiftError)
		t.Fail()
	}

	var record string
	err := Db.QueryRow("SELECT title FROM shifts WHERE title LIKE 'Test Shift'").Scan(&record)
	if err != nil {
		log.Info("Test Failed because: ", err)
		t.Fail()
	}
	if record != "Test Shift" {
		t.Fail()
	}
	cleanupShifts()
}

// TestGetShifts gets all shifts
func TestGetShiftsModel(t *testing.T) {
	ConnToDB(os.Getenv("PGURL"))

	_, delErr := Db.Query("DELETE FROM schedules WHERE title like 'Test Shift'")
	log.Info(delErr)

	spoofShifts()

	memberID := GetMemberID("testuser144@gmail.com")
	log.Info("Member ID used for testing: ", memberID)
	var scheduleID string
	schedules, getSchedErr := GetSchedules(memberID)
	if getSchedErr != nil {
		t.Error("Error getting schedule.", getSchedErr)
		t.FailNow()
	}
	scheduleID = schedules[0].ID
	log.Info("ScheduleID used to test GetShifts", scheduleID)

	s := Shift{
		Title:        "Test Shift",
		Schedule:     scheduleID,
		Start:        "715",
		End:          "1200",
		MinEnrollees: "1",
		Days:         [7]string{"", "Monday", "Tuesday", "", "", "Friday", ""},
	}

	s2 := Shift{
		Title:        "Test Shift 2",
		Schedule:     scheduleID,
		Start:        "1015",
		End:          "1700",
		MinEnrollees: "1",
		Days:         [7]string{"", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", ""},
	}

	s3 := Shift{
		Title:        "Test Shift 3",
		Schedule:     scheduleID,
		Start:        "1200",
		End:          "1200",
		MinEnrollees: "1",
		Days:         [7]string{"Sunday", "", "", "", "", "Friday", "Saturday"},
	}

	var newShiftError error
	newShiftError = CreateShift(&s)
	if newShiftError != nil {
		log.Info("Test failed while creating new shift: ", newShiftError)
		t.Fail()
	}

	newShiftError = CreateShift(&s2)
	if newShiftError != nil {
		log.Info("Test failed while creating new shift: ", newShiftError)
		t.Fail()
	}

	newShiftError = CreateShift(&s3)
	if newShiftError != nil {
		log.Info("Test failed while creating new shift: ", newShiftError)
		t.Fail()
	}

	shifts, getErr := GetShifts(scheduleID)
	if getErr != nil {
		log.Info("Failed due to: ", getErr)
		t.Fail()
	}

	if len(shifts) < 3 {
		t.Error("Expected 3 shifts and got fewer than 3.")
	}
	if len(shifts) != 3 {
		t.Error("Expected 3 shifts and did not get 3.  Got: ", len(shifts))
	}

	cleanupShifts()
}

// Helpers for shift tests
func spoofShifts() {
	m := NewMember{
		Name:      "Test Member",
		Email:     "testuser144@gmail.com",
		Email2:    "testuser144@gmail.com",
		Password:  "superduper",
		Password2: "superduper",
	}
	if CreateMember(&m) != nil {
		log.Info("Member may already be there")
	}

	l := make([]*Schedule, 1)
	l[0] = &Schedule{"", "My Morning Schedule", GetMemberID("testuser144@gmail.com")}

	for i := range l {
		if CreateSchedule(l[i]) != nil {
			log.Info("Schedule may already exist")
		}
	}
}

// cleanupShifts undoes the shift spoofing func with a cascade delete
func cleanupShifts() {
	_, err := Db.Query("DELETE FROM members WHERE email LIKE 'testuser144@gmail.com'")
	if err != nil {
		log.Info(err)
	}
}
