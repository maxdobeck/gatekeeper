package schedules

import (
	"encoding/json"
	"fmt"
	"github.com/maxdobeck/gatekeeper/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type newScheduleBody struct {
	Title    string `json:"title"`
	Owner_id string `json:"owner_id"`
}

// TestCreateNewSchedule tries to create a new schedule
func TestCreateNewSchedule(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	m := populateDb()
	// Create the schedule
	s := newScheduleBody{
		Title:    "Night Shift at Paddys",
		Owner_id: models.GetMemberID(m.Email),
	}
	b, jsonErr := json.Marshal(s)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		t.Fail()
	}
	rbody := strings.NewReader(string(b))
	req, rErr := http.NewRequest("POST", "/schedules", rbody)
	if rErr != nil {
		fmt.Println("Problem creating new schedule: ", rErr)
		t.Fail()
	}
	w := httptest.NewRecorder()
	NewSchedule(w, req)
	res := ResDetails{}
	json.Unmarshal([]byte(w.Body.String()), &res)
	var expectedMessage [1]string
	expectedMessage[0] = "Schedule created: Night Shift at Paddys"
	if res.Status != expectedMessage[0] {
		t.Error("The Schedule 'Night Shift at Paddys' was not created!")
		t.Fail()
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
		fmt.Println("Member may already be there")
	}

	l := make([]*models.Schedule, 4)
	l[0] = &models.Schedule{"", "Test Test Schedule", models.GetMemberID(m.Email)}
	l[1] = &models.Schedule{"", "My 2nd Schedule", models.GetMemberID(m.Email)}
	l[2] = &models.Schedule{"", "My 3rd Schedule", models.GetMemberID(m.Email)}
	l[3] = &models.Schedule{"", "My 4th Schedule", models.GetMemberID(m.Email)}

	for i := range l {
		if models.CreateSchedule(l[i]) != nil {
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
