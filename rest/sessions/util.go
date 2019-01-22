// Package sessions resolves session related issues
package sessions

import (
	"github.com/antonlindstrom/pgstore"
	_ "github.com/lib/pq" // github.com/lib/pq
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("super-secret-key")
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GoodSession returns true or false depending on if the session is active
func GoodSession(r *http.Request) bool {
	store, err := pgstore.NewPGStore(os.Getenv("PGURL"), key)
	check(err)
	defer store.Close()

	session, err := store.Get(r, "scheduler-session")
	check(err)

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Info("Stale session rejected: ", session)
		return false
	}
	log.Info("Session OK: ", session)
	return true
}

func CookieMemberID(r *http.Request) string {
	store, err := pgstore.NewPGStore(os.Getenv("PGURL"), key)
	check(err)
	defer store.Close()

	session, err := store.Get(r, "scheduler-session")
	check(err)

	var memberID string
	cookieData, ok := session.Values["memberID"].(string)
	if ok == true {
		memberID = cookieData
	}
	if ok != true {
		memberID = "Error"
	}
	return memberID
}
