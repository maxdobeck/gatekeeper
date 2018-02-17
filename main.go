package main

import (
	"github.com/antonlindstrom/pgstore"
	"github.com/maxdobeck/gatekeeper/authentication"
	"log"
	"net/http"
	"os"
	"time"
	"fmt"
)

func main() {
	store, err := pgstore.NewPGStore(os.Getenv("PGURL"), []byte("secret-key"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer store.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	http.HandleFunc("/validate", gatekeeper.ValidSession)
	http.HandleFunc("/login", gatekeeper.Login)
	http.HandleFunc("/logout", gatekeeper.Logout)
	
	fmt.Println("Listening on http://localhost:3030")
	http.ListenAndServe(":3030", nil)
}
