package sessions

import (
	"encoding/json"
	"github.com/gorilla/csrf"
	"github.com/maxdobeck/gatekeeper/models"
	"log"
	"net/http"
)

type resDetails struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

type curMember struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    string `json:"id"`
}

type Payload struct {
	resDetails `json:"details"`
	curMember  `json:"member"`
}

// CsrfToken will generate a CSRF Token
func CsrfToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating csrf token")
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
}

// ValidSession checks that the session is valid and can user can make requests
func ValidSession(w http.ResponseWriter, r *http.Request) {
	if GoodSession(r) != true {
		log.Println("Session is old, must log out log back in.")
		//w.WriteHeader(http.StatusUnauthorized)
		http.Error(w, "Session is expired.", http.StatusUnauthorized)
	} else {
		log.Println("Session is good.")
		w.WriteHeader(http.StatusOK)
	}
}

/* CurMember returns the currently logged in user's info for the client to consume.
It will check that the session is valid and reuturn a payload containing the member's info
based on the cookie values */
func CurMember(w http.ResponseWriter, r *http.Request) {
	if GoodSession(r) != true {
		msg := resDetails{
			Status:  "Expired session or cookie",
			Message: "Session Expired.  Log out and log back in.",
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	memberID := CookieMemberID(r)
	if memberID == "Error" {
		log.Println("No valid value in cookie.  Log out and log back in.")
		// resDetails should have the error for the client here
	}
	name := models.GetMemberName(memberID)
	email := models.GetMemberEmail(memberID)
	member := curMember{Name: name, Email: email, ID: memberID}
	log.Println("Current member based on cookie: ", member)

	msgDetails := resDetails{
		Status:  "OK",
		Message: "Member found",
	}
	msg := Payload{
		resDetails: msgDetails,
		curMember:  member,
	}
	log.Println("Payload for /CurMember: ", msg)
	json.NewEncoder(w).Encode(msg)
}
