package main

import (
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/maxdobeck/gatekeeper/authentication"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	store, err := pgstore.NewPGStore(os.Getenv("PGURL"), []byte("secret-key"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer store.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	mux := http.NewServeMux()

	mux.HandleFunc("/validate", gatekeeper.ValidSession)
	mux.HandleFunc("/login", gatekeeper.Login)
	mux.HandleFunc("/logout", gatekeeper.Logout)

	handler := cors.Default().Handler(mux)

	fmt.Println("Listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
