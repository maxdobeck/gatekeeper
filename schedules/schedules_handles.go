// Package schedules allows us to manipulate Schedules
package schedules

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

type updateSchedule struct {
	NewTitle string
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
	// Check that the owner ID matches the cookie's ID. (I.E. check that the user really is who they say they are)
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
	vars := mux.Vars(r)
	log.Println(vars["id"])

	curUser := sessions.CookieMemberID(r)
	if curUser == "Error" {
		log.Println("Problem getting member ID from cookie.  Log in and log out.")
	}
	// Check that current user is allowed to delete the schedule
	// (that the cookie session for the logged in user == the schedule owner)
	schedule, sErr := models.GetScheduleById(vars["id"])
	if sErr != nil {
		log.Println("Problem finding schedule: ", vars["id"])
		msg := ResDetails{
			Status:  "Could not find schedule.",
			Message: "Schedule does not exist.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	if curUser != schedule.OwnerID {
		msg := ResDetails{
			Status:  "Not Authorized",
			Message: "You are not the owner of this schedule.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	log.Printf("User %s is deleting schedule %s.", curUser, vars["id"])
	delErr := models.DeleteSchedule(vars["id"])
	if delErr != nil {
		msg := ResDetails{
			Status:  "Error deleting schedule",
			Message: delErr.Error(),
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	msg := ResDetails{
		Status:  "OK",
		Message: "Schedule deleted",
	}
	json.NewEncoder(w).Encode(msg)
}

func UpdateScheduleTitle(w http.ResponseWriter, r *http.Request) {
	if sessions.GoodSession(r) != true {
		msg := ResDetails{
			Status:  "Expired session or cookie",
			Message: "Session Expired.  Log out and log back in.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	// Get the variable from the url with mux
	vars := mux.Vars(r)
	if vars["id"] == "" {
		var msg ResDetails
		log.Println("Unexpected URL:", r.URL)
		msg.Status = "Error"
		msg.Message = "Bad id for schedule."
		msg.Errors = append(msg.Errors, "Path is unexpected.  Resource not found.")
		json.NewEncoder(w).Encode(msg)
		return
	}
	curUser := sessions.CookieMemberID(r)
	if curUser == "Error" {
		log.Println("Problem getting member ID from cookie.  Log in and log out.")
		msg := ResDetails{
			Status:  "Expired session or cookie",
			Message: "Problem with session.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	// Check that current user is allowed to touch the schedule
	// (that the cookie session for the logged in user == the schedule owner)
	schedule, sErr := models.GetScheduleById(vars["id"])
	if sErr != nil {
		log.Println("Problem finding schedule: ", vars["id"])
		msg := ResDetails{
			Status:  "Could not find schedule.",
			Message: "Schedule does not exist.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	if curUser != schedule.OwnerID {
		msg := ResDetails{
			Status:  "Not Authorized",
			Message: "You are not the owner of this schedule.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	var titleUpdate updateSchedule
	err := json.NewDecoder(r.Body).Decode(&titleUpdate)
	if err != nil {
		log.Println("Error decoding body >>", err)
		msg := ResDetails{
			Status:  "Error.",
			Message: "Error decoding body.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	log.Printf("User %s is updating schedule %s with new title: %s", curUser, vars["id"], titleUpdate.NewTitle)
	updateErr := models.UpdateScheduleTitle(vars["id"], titleUpdate.NewTitle)
	if updateErr != nil {
		msg := ResDetails{
			Status:  "Error changing schedule title.",
			Message: updateErr.Error(),
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	msg := ResDetails{
		Status:  "OK",
		Message: fmt.Sprintf("Title Updated: %s", titleUpdate.NewTitle),
	}
	json.NewEncoder(w).Encode(msg)
}

// Find Schedule based on the specified schedule ID
func FindScheduleByID(w http.ResponseWriter, r *http.Request) {
	if sessions.GoodSession(r) != true {
		msg := ResDetails{
			Status:  "Expired session or cookie",
			Message: "Session Expired.  Log out and log back in.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	// Get the variable from the url with mux
	vars := mux.Vars(r)
	if vars["id"] == "" {
		var msg ResDetails
		log.Println("Unexpected URL:", r.URL)
		msg.Status = "Error"
		msg.Message = "Bad id for schedule."
		msg.Errors = append(msg.Errors, "Path is unexpected.  Resource not found.")
		json.NewEncoder(w).Encode(msg)
		return
	}
	schedule, sErr := models.GetScheduleById(vars["id"])
	if sErr != nil {
		log.Println("Problem finding schedule: ", vars["id"])
		msg := ResDetails{
			Status:  "Could not find schedule.",
			Message: "Schedule does not exist or could not be found.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	msg := Payload{}
	msg.FoundSchedules = append(msg.FoundSchedules, schedule)
	details := ResDetails{
		Status:  "OK",
		Message: "Schedule Found",
	}
	msg.ResDetails = details
	json.NewEncoder(w).Encode(msg)
}

// Find All Schedules based on member ID
func FindSchedulesByOwner(w http.ResponseWriter, r *http.Request) {
	if sessions.GoodSession(r) != true {
		msg := ResDetails{
			Status:  "Expired session or cookie",
			Message: "Session Expired.  Log out and log back in.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	// Get the variable from the url with mux
	vars := mux.Vars(r)
	if vars["id"] == "" {
		var msg ResDetails
		log.Println("Unexpected URL:", r.URL)
		msg.Status = "Error"
		msg.Message = "Bad id for member."
		msg.Errors = append(msg.Errors, "Path is unexpected.  Resource not found.")
		json.NewEncoder(w).Encode(msg)
		return
	}
	schedules, sErr := models.GetSchedules(vars["id"])
	if sErr != nil {
		log.Println("Problem finding schedule: ", vars["id"])
		msg := ResDetails{
			Status:  "Could not find schedule.",
			Message: "Schedule does not exist or could not be found.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	msg := Payload{}
	msg.FoundSchedules = schedules
	log.Println("Schedules found: ", schedules)
	details := ResDetails{
		Status:  "OK",
		Message: "Schedule Found",
	}
	msg.ResDetails = details
	log.Println("Payload from Finding Schedule by Onwer: ", msg)
	json.NewEncoder(w).Encode(msg)
}
