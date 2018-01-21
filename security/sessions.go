package gatekeeper

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// ValidSession checks if the session is authenticated and still active
func ValidSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "scheduler-session")
	if err != nil {
		panic(err)
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Is this session valid: false", http.StatusUnauthorized)
		return
	}

	// Return message
	fmt.Fprintln(w, "Is this session valid: true")
}

// Login gets a new session for the user if the credential check passes
func Login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "scheduler-session")
	if err != nil {
		panic(err)
	}
	// Authenticate based on incoming http request
	if passwordsMatch(r) != true {
		http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}
	
	// Assume passwords matched, add a function to get the memberID based ont the email
	fmt.Println(getMemberID("Jack_Reinger@yahoo.com"))
	
	// Respond with the proper content type and the memberID
	w.Header().Set("Content-Type", "text/plain") // TODO convert this to application/json
	w.WriteHeader(http.StatusOK)
	memberID := "1234"
	fmt.Fprintf(w, memberID)
	// w.Write() // Alternative to fprintf.  Needs []byte of marshalled JSON
	
	// Set cookie values and save
	session.Values["authenticated"] = true
	session.Save(r, w)
}

// Logout destroys the session
func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "scheduler-session")
	if err != nil {
		panic(err)
	}

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
