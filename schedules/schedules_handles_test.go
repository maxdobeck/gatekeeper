package schedules

import (
	"encoding/json"
	"fmt"
	"github.com/maxdobeck/gatekeeper/authentication"
	"github.com/maxdobeck/gatekeeper/models"
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
		fmt.Println(jsonErr)
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
		fmt.Println("Problem creating new request: ", rErr)
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
	if res.Status != expectedMessage[0] {
		fmt.Println("Response Status: ", res.Status)
		fmt.Printf("The Schedule '%s' was not created!\n", s.Title)
		t.Fail()
	}
	cleanupDb()
}

// Update the specified Schedule's Title
/*func TestUpdateScheduleTitle(t *testing.T) {

}

// Delete the specified schedule
func TestDeleteSchedule(t *testing.T) {

}

// Find all schedules owned by the specified member
func TestGetScheduleByOwner(t *testing.T) {

}

// Find a schedule based on the specified ID
func TestGetScheduleByID(t *testing.T) {

}

// Try and find a schedule that doesn't exist ensure proper error is returned
func TestGetNonexistentScheduleByID(t *test.T) {

}
*/

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
	l[0] = &models.Schedule{"", "Test Test Schedule", models.GetMemberID(m.Email)}
	l[1] = &models.Schedule{"", "My 2nd Schedule", models.GetMemberID(m.Email)}
	l[2] = &models.Schedule{"", "My 3rd Schedule", models.GetMemberID(m.Email)}
	l[3] = &models.Schedule{"", "My 4th Schedule", models.GetMemberID(m.Email)}

	for i := range l {
		err := models.CreateSchedule(l[i])
		if err != nil {
			fmt.Println("Schedule may already exist")
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
