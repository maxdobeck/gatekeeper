package main

import (
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/maxdobeck/gatekeeper/authentication"
	"github.com/maxdobeck/gatekeeper/members"
	"github.com/maxdobeck/gatekeeper/models"
	"github.com/maxdobeck/gatekeeper/schedules"
	"github.com/maxdobeck/gatekeeper/sessions"
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

	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)

	var allowedDomains []string // need env variables for this
	if os.Getenv("GO_ENV") == "dev" {
		allowedDomains = []string{"127.0.0.1:3000", "http://localhost:3000", "127.0.0.1:3000", "127.0.0.1:3050"}
	} else if os.Getenv("GO_ENV") == "test" {
		allowedDomains = []string{"https://s3-sih-test.s3-website-us-west-1.amazonaws.com"}
	} else if os.Getenv("GO_ENV") == "prod" {
		allowedDomains = []string{"https://schedulingishard.com", "https://www.schedulingishard.com"}
	}

	CSRF := csrf.Protect(
		[]byte("32-byte-long-auth-key"),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("scheduler_csrf"),
		csrf.Secure(false), // Disabled for localhost non-https debugging
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedDomains,
		AllowedMethods:   []string{"PUT", "POST", "PATCH", "DELETE", "GET"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-CSRF-Token"},
		ExposedHeaders:   []string{"X-CSRF-Token"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	r := mux.NewRouter()
	// Authentication Routes
	r.HandleFunc("/csrftoken", sessions.CsrfToken).Methods("GET")
	r.HandleFunc("/login", authentication.Login).Methods("POST")
	r.HandleFunc("/logout", authentication.Logout).Methods("DELETE")
	// Session Routes
	r.HandleFunc("/validsession", sessions.ValidSession).Methods("GET")
	r.HandleFunc("/curmember", sessions.CurMember).Methods("GET")
	// Member CRUD Routes
	r.HandleFunc("/members", members.SignupMember).Methods("POST")
	r.HandleFunc("/members/{id}/email", members.UpdateMemberEmail).Methods("PUT")
	r.HandleFunc("/members/{id}/name", members.UpdateMemberName).Methods("PUT")
	// r.HandleFunc("/members/{id}", members.DeleteMember).Methods("DELETE")
	// Schedules CRUD Routes
	r.HandleFunc("/schedules", schedules.NewSchedule).Methods("POST")
	r.HandleFunc("/schedules/{id}", schedules.FindScheduleByID).Methods("GET")
	r.HandleFunc("/schedules/owner/{id}", schedules.FindSchedulesByOwner).Methods("GET")
	r.HandleFunc("/schedules/{id}/title", schedules.UpdateScheduleTitle).Methods("PATCH")
	r.HandleFunc("/schedules/{id}", schedules.DeleteScheduleByID).Methods("DELETE")
	// Shifts Routes
	//  r.HandleFunc("/schedules/{scheduleid}/shifts", shifts.NewShift).Methods("POST")
	//  r.HandleFunc("/schedules/{scheduleid}/shifts", shifts.GetShifts).Methods("GET")
	//  r.HandleFunc("/schedules/{scheduleid}/shifts/{shiftid}", shifts.Get).Methods("GET")
	//  r.HandleFunc("/schedules/{scheduleid}/shifts/{shiftid}", shifts.Delete).Methods("DELETE")
	//  r.HandleFunc("/schedules/{scheduleid}/shifts/{shiftid}", shifts.Update).Methods("PATCH")
	// Middleware
	n := negroni.Classic()
	n.Use(c)
	n.UseHandler(CSRF(r))

	var hostURL string // need env variables here too
	if os.Getenv("GO_ENV") == "prod" {
		hostURL = "https://rugged-wind-cave-81042.herokuapp.com"
	} else if os.Getenv("GO_ENV") == "test" {
		hostURL = "https://shielded-stream-75107.herokuapp.com/"
	} else if os.Getenv("GO_ENV") == "dev" {
		hostURL = "http://localhost"
	}
	port := os.Getenv("PORT")
	log.Println("Listening on: ", hostURL+":"+port)
	log.Fatal(http.ListenAndServe(":"+port, context.ClearHandler(n)))
}
