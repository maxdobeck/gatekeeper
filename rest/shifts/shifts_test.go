package shifts

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxdobeck/gatekeeper/models"
	"github.com/maxdobeck/gatekeeper/rest/authentication"
	"github.com/maxdobeck/gatekeeper/rest/schedules"
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
		fmt.Println("Problem creating new request: ", rErr)
		t.Fail()
	}
	// Add the cookie from the newly created session to the request
	req.AddCookie(wLogin.Result().Cookies()[0])

	// Setup a router and test the handle
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/schedules/owner/{id}", schedules.FindSchedulesByOwner)
	router.ServeHTTP(w, req)

	res := schedules.Payload{}
	json.Unmarshal([]byte(w.Body.String()), &res)
	if len(res.FoundSchedules) < 4 {
		t.Errorf("Actual res: %s", res)
		t.Error("Not all four schedules were found.")
		t.Fail()
	}
	fmt.Println("Schedule ID we'll be using: ", res.FoundSchedules[0].ID)
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
		fmt.Println("Member may already be there")
	}

	l := make([]*models.Schedule, 4)
	l[0] = &models.Schedule{ID: "", Title: "Test Test Schedule", OwnerID: models.GetMemberID(m.Email)}
	l[1] = &models.Schedule{ID: "", Title: "My 2nd Schedule", OwnerID: models.GetMemberID(m.Email)}
	l[2] = &models.Schedule{ID: "", Title: "My 3rd Schedule", OwnerID: models.GetMemberID(m.Email)}
	l[3] = &models.Schedule{ID: "", Title: "My 4th Schedule", OwnerID: models.GetMemberID(m.Email)}

	for i := range l {
		err := models.CreateSchedule(l[i])
		if err != nil {
			fmt.Println("Schedule may already exist and you should be able to ignore any errors about duplicate keys.")
		}
	}
	return m
}

// cleanupDb undoes the populateDb
func cleanupDb() {
	_, err := models.Db.Query("DELETE FROM members WHERE email LIKE 'frank@paddys.com'")
	if err != nil {
		fmt.Println(err)
	}
}
