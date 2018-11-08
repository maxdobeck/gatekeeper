package shifts

import (
	"encoding/json"
	"fmt"
	"github.com/maxdobeck/gatekeeper/models"
	"github.com/maxdobeck/gatekeeper/rest"
	"github.com/maxdobeck/gatekeeper/rest/sessions"
	"log"
	"net/http"
)

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
		log.Println("Error decoding body >>", err)
		msg := rest.ResDetails{
			Status:  "Error",
			Message: "Couldn't decode schedule",
			Errors:  []string{"Problem decoding"},
		}
		log.Println(msg)
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
		log.Println(msg)
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

	// var s models.Shift
	// var err error

	// msg := rest.ResDetails{
	// 	Status:  "OK",
	// 	Message: fmt.Sprintf("Shifts found: %s", shifts),
	// }
	// json.NewEncoder(w).Encode(msg)
}
