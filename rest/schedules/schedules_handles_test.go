package schedules

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxdobeck/gatekeeper/models"
	"github.com/maxdobeck/gatekeeper/rest/authentication"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type newScheduleBody struct {
	Title   string `json:"Title"`
	OwnerID string `json:"OwnerID"`
}

// TestCreateNewSchedule tries to create a new schedule
func TestCreateNewSchedule(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	m := populateDb()
	// Create the schedule
	s := newScheduleBody{
		Title:   "Night Shift at Paddys",
		OwnerID: models.GetMemberID(m.Email),
	}
	b, jsonErr := json.Marshal(s)
	if jsonErr != nil {
		log.Info(jsonErr)
		t.Fail()
	}

	// Login to start a session
	loginBody := strings.NewReader(`{"email": "frank@paddys.com", "password": "superduper"}`)
	loginReq, loginErr := http.NewRequest("POST", "/login", loginBody)
	if loginErr != nil {
		t.Fail()
	}
	wLogin := httptest.NewRecorder()
	authentication.Login(wLogin, loginReq)

	rbody := strings.NewReader(string(b))
	req, rErr := http.NewRequest("POST", "/schedules", rbody)
	if rErr != nil {
		log.Info("Problem creating new request: ", rErr)
		t.Fail()
	}
	// Add the cookie from the newly created session to the request
	req.AddCookie(wLogin.Result().Cookies()[0])
	w := httptest.NewRecorder()
	NewSchedule(w, req)
	res := ResDetails{}
	json.Unmarshal([]byte(w.Body.String()), &res)
	var expectedMessage [1]string
	expectedMessage[0] = fmt.Sprintf("Schedule created: %s", s.Title)
	if res.Message != expectedMessage[0] {
		log.Info("Response Status: ", res.Status)
		log.Warnf("The Schedule '%s' was not created!\n", s.Title)
		t.Fail()
	}
	cleanupDb()
}

// Update the specified Schedule's Title
func TestUpdateScheduleTitle(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	m := populateDb()
	ownerID := models.GetMemberID(m.Email)
	var scheduleID string

	findErr := models.Db.QueryRow("SELECT id FROM schedules WHERE owner_id = $1 LIMIT 1", ownerID).Scan(&scheduleID)
	if findErr != nil {
		t.Errorf("The shedule %s could not be found ", scheduleID)
	}

	// Login to start a session
	loginBody := strings.NewReader(`{"email": "frank@paddys.com", "password": "superduper"}`)
	loginReq, loginErr := http.NewRequest("POST", "/login", loginBody)
	if loginErr != nil {
		t.Fail()
	}
	wLogin := httptest.NewRecorder()
	authentication.Login(wLogin, loginReq)
	// Build the request to test
	body := strings.NewReader(`{"newtitle": "New Schedule Title"}`)
	req, rErr := http.NewRequest("PATCH", "/schedules/"+scheduleID+"/title", body)
	if rErr != nil {
		log.Info("Problem creating new request: ", rErr)
		t.Fail()
	}
	// Add the cookie from the newly created session to the request
	req.AddCookie(wLogin.Result().Cookies()[0])
	// Setup a router and test the handle
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/{id}/title", UpdateScheduleTitle)
	router.ServeHTTP(w, req)

	res := ResDetails{}
	json.Unmarshal([]byte(w.Body.String()), &res)
	if res.Status != "OK" {
		t.Error("Error updating Schedule Title")
		t.Fail()
	}
	if res.Message != "Title Updated: New Schedule Title" {
		t.Errorf("Actual Response: %s", res)
		t.Errorf("The schedule: %s title was not updated.", scheduleID)
		t.Fail()
	}
	cleanupDb()
}

// Delete the specified schedule
func TestDeleteSchedule(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	m := populateDb()
	ownerID := models.GetMemberID(m.Email)
	var scheduleID string

	findErr := models.Db.QueryRow("SELECT id FROM schedules WHERE owner_id = $1 LIMIT 1", ownerID).Scan(&scheduleID)
	if findErr != nil {
		t.Errorf("The schedule %s could not be found ", scheduleID)
	}

	// Login to start a session
	loginBody := strings.NewReader(`{"email": "frank@paddys.com", "password": "superduper"}`)
	loginReq, loginErr := http.NewRequest("POST", "/login", loginBody)
	if loginErr != nil {
		t.Fail()
	}
	wLogin := httptest.NewRecorder()
	authentication.Login(wLogin, loginReq)
	// Build the request to test
	req, rErr := http.NewRequest("DELETE", "/schedules/"+scheduleID, nil)
	if rErr != nil {
		log.Info("Problem creating new request: ", rErr)
		t.Fail()
	}
	// Add the cookie from the newly created session to the request
	req.AddCookie(wLogin.Result().Cookies()[0])
	// Setup a router and test the handle
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/{id}", DeleteScheduleByID)
	router.ServeHTTP(w, req)

	res := Payload{}
	json.Unmarshal([]byte(w.Body.String()), &res)
	if res.ResDetails.Message != "Schedule deleted" {
		t.Errorf("The schedule %s could not be deleted ", scheduleID)
		t.Fail()
	}
	cleanupDb()
}

// Find all schedules owned by the specified member
func TestFindSchedulesByOwner(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	m := populateDb()

	// Login to start a session
	loginBody := strings.NewReader(`{"email": "frank@paddys.com", "password": "superduper"}`)
	loginReq, loginErr := http.NewRequest("POST", "/login", loginBody)
	if loginErr != nil {
		t.Fail()
	}
	wLogin := httptest.NewRecorder()
	authentication.Login(wLogin, loginReq)
	memberID := models.GetMemberID(m.Email)
	req, rErr := http.NewRequest("GET", "/schedules/owner/"+memberID, nil)
	if rErr != nil {
		log.Info("Problem creating new request: ", rErr)
		t.Fail()
	}
	// Add the cookie from the newly created session to the request
	req.AddCookie(wLogin.Result().Cookies()[0])

	// Setup a router and test the handle
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/owner/{id}", FindSchedulesByOwner)
	router.ServeHTTP(w, req)

	res := Payload{}
	json.Unmarshal([]byte(w.Body.String()), &res)
	if len(res.FoundSchedules) < 4 {
		t.Errorf("Actual res: %s", res)
		t.Error("Not all four schedules were found.")
		t.Fail()
	}
	cleanupDb()
}

// Find a schedule based on the specified ID
func TestFindScheduleByID(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	m := populateDb()
	ownerID := models.GetMemberID(m.Email)
	var targetID string

	findErr := models.Db.QueryRow("SELECT id FROM schedules WHERE owner_id = $1 LIMIT 1", ownerID).Scan(&targetID)
	if findErr != nil {
		t.Error("Problem finding a schedule in the database.")
	}
	// Login to grab a valid session cookie
	loginBody := strings.NewReader(`{"email": "frank@paddys.com", "password": "superduper"}`)
	loginReq, loginErr := http.NewRequest("POST", "/login", loginBody)
	if loginErr != nil {
		t.Fail()
	}
	wLogin := httptest.NewRecorder()
	authentication.Login(wLogin, loginReq)
	req, rErr := http.NewRequest("GET", "/schedules/"+targetID, nil)
	if rErr != nil {
		log.Info("Problem creating new request: ", rErr)
		t.Fail()
	}
	// Add the cookie from the newly created session to the request
	req.AddCookie(wLogin.Result().Cookies()[0])

	// Setup a router and test the handle
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/{id}", FindScheduleByID)
	router.ServeHTTP(w, req)

	res := SinglePayload{}
	json.Unmarshal([]byte(w.Body.String()), &res)
	if res.FoundSchedule.ID != targetID {
		t.Error("Bad schedule was returned in the payload")
		t.Fail()
	}
	if res.ResDetails.Status != "OK" {
		t.Errorf("The schedule %s could not be found", targetID)
	}

	cleanupDb()
}

// Try and find a schedule that doesn't exist ensure proper error is returned
func TestGetNonexistentScheduleByID(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	m := populateDb()
	var targetID = "17wrong-id-type"
	// Login to grab a valid session cookie
	loginBody := strings.NewReader(`{"email": "frank@paddys.com", "password": "superduper"}`)
	loginReq, loginErr := http.NewRequest("POST", "/login", loginBody)
	if loginErr != nil {
		t.Fail()
	}
	wLogin := httptest.NewRecorder()
	authentication.Login(wLogin, loginReq)
	req, rErr := http.NewRequest("GET", "/schedules/owners/"+models.GetMemberID(m.Email), nil)
	if rErr != nil {
		log.Info("Problem creating new request: ", rErr)
		t.Fail()
	}
	// Add the cookie from the newly created session to the request
	req.AddCookie(wLogin.Result().Cookies()[0])

	// Setup a router and test the handle
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/{id}", FindScheduleByID)
	router.ServeHTTP(w, req)

	res := Payload{}
	json.Unmarshal([]byte(w.Body.String()), &res)
	if len(res.FoundSchedules) > 0 {
		t.Errorf("A schedule was somehow returned when nothing was expected.  %s", res)
		t.Fail()
	}
	if res.ResDetails.Status == "schedule found" {
		t.Error("Response: ", res)
		t.Errorf("The schedule %s was found when it doesn't exist", targetID)
	}

	cleanupDb()
}

// Helpers
func populateDb() models.NewMember {
	m := models.NewMember{
		Name:      "Frank",
		Email:     "frank@paddys.com",
		Email2:    "frank@paddys.com",
		Password:  "superduper",
		Password2: "superduper",
	}
	if models.CreateMember(&m) != nil {
		log.Info("Member may already be there")
	}

	l := make([]*models.Schedule, 4)
	l[0] = &models.Schedule{ID: "", Title: "Test Test Schedule", OwnerID: models.GetMemberID(m.Email)}
	l[1] = &models.Schedule{ID: "", Title: "My 2nd Schedule", OwnerID: models.GetMemberID(m.Email)}
	l[2] = &models.Schedule{ID: "", Title: "My 3rd Schedule", OwnerID: models.GetMemberID(m.Email)}
	l[3] = &models.Schedule{ID: "", Title: "My 4th Schedule", OwnerID: models.GetMemberID(m.Email)}

	for i := range l {
		err := models.CreateSchedule(l[i])
		if err != nil {
			log.Info("Schedule may already exist and you should be able to ignore any errors about duplicate keys.")
		}
	}
	return m
}

// cleanupDb undoes the populateDb
func cleanupDb() {
	_, err := models.Db.Query("DELETE FROM members WHERE email LIKE 'frank@paddys.com'")
	if err != nil {
		log.Info(err)
	}
}
