package shifts

import (
	"encoding/json"
	"github.com/maxdobeck/gatekeeper/rest"
	"github.com/maxdobeck/gatekeeper/sessions"
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
	log.Println("Creating new shift: ")
}
