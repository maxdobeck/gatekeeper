package shifts

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxdobeck/gatekeeper/models"
	"github.com/maxdobeck/gatekeeper/rest"
	"github.com/maxdobeck/gatekeeper/rest/sessions"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Payload contains the client facing messages to an action
type Payload struct {
	rest.ResDetails
	FoundShifts []models.ShiftPayload
}

// New attempts to make a new shift for a schedule
func New(w http.ResponseWriter, r *http.Request) {
	if sessions.GoodSession(r) != true {
		msg := rest.ResDetails{
			Status:  "Expired session or cookie",
			Message: "Session Expired.  Log out and log back in.",
			Errors:  []string{"Session Expired"},
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	var s models.Shift
	var err error
	err = json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Info("Error decoding body >>", err, r.Body)
		msg := rest.ResDetails{
			Status:  "Error",
			Message: "Couldn't decode schedule",
			Errors:  []string{"Problem decoding"},
		}
		log.Info(msg)
		json.NewEncoder(w).Encode(msg)
		return
	}
	err = models.CreateShift(&s)
	if err != nil {
		msg := rest.ResDetails{
			Status:  "Error",
			Message: fmt.Sprintf("Couldn't create schedule in database: %s", s),
			Errors:  []string{"Problem creating record", err.Error()},
		}
		log.Info(msg)
		json.NewEncoder(w).Encode(msg)
		return
	}
	msg := rest.ResDetails{
		Status:  "OK",
		Message: fmt.Sprintf("Shift created: %s", s.Title),
	}
	json.NewEncoder(w).Encode(msg)
}

// FindAll will attempt to find all shifts based on the ScheduleID
func FindAll(w http.ResponseWriter, r *http.Request) {
	if sessions.GoodSession(r) != true {
		msg := rest.ResDetails{
			Status:  "Expired session or cookie",
			Message: "Session Expired.  Log out and log back in.",
			Errors:  []string{"Session Expired"},
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	// Get the scheduleID from the body and then find all shifts
	log.Info(r.Body)

	vars := mux.Vars(r)
	log.Info("Get All Shifts for this Schedule: ", vars["scheduleid"], vars)

	curUser := sessions.CookieMemberID(r)
	if curUser == "Error" {
		log.Info("Problem getting member ID from cookie.  Log in and log out.")
		msg := rest.ResDetails{
			Status:  "Expired session or cookie",
			Message: "Session Expired.  Log out and log back in.",
			Errors:  []string{"Session Expired"},
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	foundShifts, getAllErr := models.GetShifts(vars["scheduleid"])
	if getAllErr != nil {
		msg := rest.ResDetails{
			Status:  "Error",
			Message: "Problem Getting All Shifts for Schedule",
			Errors:  []string{getAllErr.Error()},
		}
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Payload{}
	msg.FoundShifts = foundShifts
	details := rest.ResDetails{
		Status:  "OK",
		Message: fmt.Sprintf("%v Shift(s) found", len(foundShifts)),
	}
	msg.ResDetails = details
	json.NewEncoder(w).Encode(msg)
}
