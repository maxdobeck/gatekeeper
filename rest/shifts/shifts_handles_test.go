package shifts

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/maxdobeck/gatekeeper/models"
	"github.com/maxdobeck/gatekeeper/rest/authentication"
	"github.com/maxdobeck/gatekeeper/rest/schedules"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestFindAllShifts uses a scheduleID to get all the shifts for that schedule
func TestFindAllShifts(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	m := populateDb()

	// Login to start a session
	loginBody := strings.NewReader(`{"email": "frankshift@paddys.com", "password": "superduper"}`)
	loginReq, loginErr := http.NewRequest("POST", "/login", loginBody)
	if loginErr != nil {
		t.Fail()
	}
	wLogin := httptest.NewRecorder()
	authentication.Login(wLogin, loginReq)
	memberID := models.GetMemberID(m.Email)
	log.Info("Member we'll be using: ", memberID, m.Email)
	req, rErr := http.NewRequest("GET", "/schedules/owner/"+memberID, nil)
	if rErr != nil {
		log.Fatal("Problem creating new request: ", rErr)
		t.Fail()
	}
	// Add the cookie from the newly created session to the request
	req.AddCookie(wLogin.Result().Cookies()[0])

	// Setup a router and test the handle
	scheduleRecorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/owner/{id}", schedules.FindSchedulesByOwner)
	router.ServeHTTP(scheduleRecorder, req)

	schedRes := schedules.Payload{}
	log.Info("RAW Schedule Recorder Body", scheduleRecorder.Body.String())
	json.Unmarshal([]byte(scheduleRecorder.Body.String()), &schedRes)
	log.Info("SCHEDULE RESP JSON Payload: ", schedRes)
	if schedRes.ResDetails.Status != "OK" {
		t.Errorf("Actual /schedules/owner/{id} response: %s to this request %s", scheduleRecorder.Body.String(), req.URL)
		t.Error("Not all schedules were found.", schedRes)
		t.Fail()
	}
	log.Info("Schedule ID we'll be using: ", schedRes.FoundSchedules[0].ID)

	shiftReq, shiftReqErr := http.NewRequest("GET", "/schedules/"+schedRes.FoundSchedules[0].ID+"/shifts", nil)
	if shiftReqErr != nil {
		log.Fatal("Problem creating new request: ", shiftReqErr)
		t.Fail()
	}
	// Add the cookie from the newly created session to the request
	shiftReq.AddCookie(wLogin.Result().Cookies()[0])

	shiftRecorder := httptest.NewRecorder()
	shiftRouter := mux.NewRouter()
	shiftRouter.HandleFunc("/schedules/{scheduleid}/shifts", FindAll)
	shiftRouter.ServeHTTP(shiftRecorder, shiftReq)

	// Actual test: We're looking for 1 shift
	log.Info("schedules/id/shifts output: ", shiftRecorder.Body.String())
	shiftPayload := Payload{}
	json.Unmarshal([]byte(shiftRecorder.Body.String()), &shiftPayload)
	if len(shiftPayload.FoundShifts) != 1 {
		t.Error("Did not get our 1 target shift.", shiftPayload)
	}
	cleanupDb()
}

// Helpers
func populateDb() models.NewMember {
	m := models.NewMember{
		Name:      "Frank",
		Email:     "frankshift@paddys.com",
		Email2:    "frankshift@paddys.com",
		Password:  "superduper",
		Password2: "superduper",
	}
	if models.CreateMember(&m) != nil {
		log.Fatal("Member may already be there")
	}

	l := make([]*models.Schedule, 1)
	l[0] = &models.Schedule{ID: "", Title: "Test Test Schedule", OwnerID: models.GetMemberID(m.Email)}
	// l[1] = &models.Schedule{ID: "", Title: "My 2nd Schedule", OwnerID: models.GetMemberID(m.Email)}
	// l[2] = &models.Schedule{ID: "", Title: "My 3rd Schedule", OwnerID: models.GetMemberID(m.Email)}
	// l[3] = &models.Schedule{ID: "", Title: "My 4th Schedule", OwnerID: models.GetMemberID(m.Email)}

	for i := range l {
		err := models.CreateSchedule(l[i])
		if err != nil {
			log.Fatal("Schedule may already exist and you should be able to ignore any errors about duplicate keys.")
		}
	}

	franksSchedules, getErr := models.GetSchedules(models.GetMemberID(m.Email))
	if getErr != nil {
		log.Info("problem getting all of Franks Schedules!", getErr)
	}

	target := models.Shift{
		Title:        "Morning Shift",
		Start:        "05:00",
		End:          "10:00",
		Stop:         "2099-01-01",
		MinEnrollees: "1",
		Schedule:     franksSchedules[0].ID,
		Days:         [7]string{"", "monday", "tuesday", "wednesday", "thursday", "friday", ""},
	}
	shiftErr := models.CreateShift(&target)
	if shiftErr != nil {
		log.Fatal("Target Shift may already exist but this could be a real error!", shiftErr)
	}
	shifts, _ := models.GetShifts(franksSchedules[0].ID)

	log.Info("All of Franks first schedule shifts: ", shifts)

	return m
}

// cleanupDb undoes the populateDb
func cleanupDb() {
	_, err := models.Db.Query("DELETE FROM members WHERE email LIKE 'frankshift@paddys.com'")
	if err != nil {
		log.Fatal(err)
	}
}
