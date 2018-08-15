// Package schedules allows us to manipulate Schedules
package schedules

import (
	"encoding/json"
	"fmt"
	"github.com/maxdobeck/gatekeeper/models"
	"github.com/maxdobeck/gatekeeper/sessions"
	"log"
	"net/http"
)

// ResDetails contains the response status, messages, and any errors
type ResDetails struct {
	Status  string
	Message string
	Errors  []string
}

type Payload struct {
	ResDetails
	FoundSchedules []models.Schedule
}

// NewSchedule is used to make a new schedule
func NewSchedule(w http.ResponseWriter, r *http.Request) {
	var newScheduleErrors []string
	if sessions.GoodSession(r) != true {
		msg := ResDetails{
			Status:  "Expired session or cookie",
			Message: "Session Expired.  Log out and log back in.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	var s models.Schedule
	// var errors []string
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Println("Problem decoding incoming Schedule", err)
	}
	scheduleErr := models.CreateSchedule(&s)
	if scheduleErr != nil {
		log.Println("Problem making schedule: ", scheduleErr, s)
		newScheduleErrors = append(newScheduleErrors, fmt.Sprintf("Error creating schedule %s", scheduleErr))
		msg := ResDetails{
			Status:  fmt.Sprintf("Problem creating schedule: %s", s.Title),
			Message: fmt.Sprintf("Error: %s", scheduleErr),
			Errors:  newScheduleErrors,
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	msg := ResDetails{
		Status: fmt.Sprintf("Schedule created: %s", s.Title),
	}
	json.NewEncoder(w).Encode(msg)
}

// Delete schedule by specified ID if owner made request
func DeleteScheduleByID(w http.ResponseWriter, r *http.Request) {
	if sessions.GoodSession(r) != true {
		msg := ResDetails{
			Status:  "Expired session or cookie",
			Message: "Session Expired.  Log out and log back in.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	// Check that current user is allowed to delete the schedule (They are the owner)
	curUser := sessions.CookieMemberID(r)
	if curUser == "Error" {
		log.Println("Problem getting member ID from cookie.  Log in and log out.")
	}
	log.Println("Current user ID: ", curUser)
	msg := ResDetails{
		Status:  "OK",
		Message: "Schedule deleted",
	}
	json.NewEncoder(w).Encode(msg)
}

func UpdateScheduleTitle(w http.ResponseWriter, r *http.Request) {

}

// Find Schedule based on the specified schedule ID
func FindScheduleByID(w http.ResponseWriter, r *http.Request) {

}

// Find All Schedules based on member ID
func FindScheduleByOwner(w http.ResponseWriter, r *http.Request) {

}
