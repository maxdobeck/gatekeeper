// Package schedules allows us to manipulate Schedules
package schedules

import (
	// "github.com/maxdobeck/gatekeeper/models"
	"encoding/json"
	"github.com/maxdobeck/gatekeeper/sessions"
	"net/http"
)

type resDetails struct {
	Status  string
	Message []string
}

// CreateSchedule is used to make a new schedule
func CreateSchedule(w http.ResponseWriter, r *http.Request) {
	if sessions.GoodSession(r) != true {
		msg := resDetails{
			Status:  "Expired session or cookie",
			Message: []string{"Session Expired.  Log out and log back in."},
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
}
