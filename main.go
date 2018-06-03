package main

import (
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/maxdobeck/gatekeeper/authentication"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
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

	CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:3030", "http://localhost:3030", "https://schedulingishard.com", "https://www.schedulingishard.com"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/validate", gatekeeper.ValidSession)
	mux.HandleFunc("/login", gatekeeper.Login)
	mux.HandleFunc("/logout", gatekeeper.Logout)

	n := negroni.Classic()
	n.Use(c)
	n.UseHandler(CSRF(mux))

	fmt.Println("Listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", context.ClearHandler(n)))
}
