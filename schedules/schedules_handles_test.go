package schedules

import (
	"fmt"
	"github.com/maxdobeck/gatekeeper/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestCreateNewSchedule tries to create a new schedule
func TestCreateNewSchedule(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	populateDb()
	// Create the schedule
	rbody := strings.NewReader(`{"title": "Night Shift at Paddys", "owner_id": }`)
	req, rErr := http.NewRequest("POST", "/schedules", rbody)
	if rErr != nil {
		fmt.Println("Problem creating new schedule: ", rErr)
		t.Fail()
	}
	w := httptest.NewRecorder()
	NewSchedule(w, req)

	cleanupDb()
	/*
		// Signup a user
		signupBody := strings.NewReader(`{"email": "testValidCreds@gmail.com", "email2":"testValidCreds@gmail.com", "password": "supersecret", "password2":"supersecret", "name":"Valid User Signup"}`)
		signupReq, signupErr := http.NewRequest("POST", "/members", signupBody)
		if signupErr != nil {
			t.Fail()
		}
		wSignup := httptest.NewRecorder()
		SignupMember(wSignup, signupReq)

		fmt.Println("Now try and sign up user with same email")
		dupBody := strings.NewReader(`{"email": "testValidCreds@gmail.com", "email2":"testValidCreds@gmail.com", "password": "supersecret", "password2":"supersecret", "name":"Valid User Signup"}`)
		dupReq, dupSignupErr := http.NewRequest("POST", "/members", dupBody)
		if dupSignupErr != nil {
			t.Fail()
		}
		//Signup same member again
		wSignup2 := httptest.NewRecorder()
		SignupMember(wSignup2, dupReq)
		actualRes := memberSignup{}
		json.Unmarshal([]byte(wSignup2.Body.String()), &actualRes)
		var expectedRes [1]string
		expectedRes[0] = "Email is already in use."
		if actualRes.Errors[0] != expectedRes[0] {
			t.Error("SignupMember allowed a duplicate email through.", wSignup2.Body)
		} */
}

// Helpers
func populateDb() {
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
}

// cleanupDb undoes the populateDb
func cleanupDb() {
	_, err := models.Db.Query("DELETE FROM members WHERE email LIKE 'frank@paddys.com'")
	if err != nil {
		fmt.Println(err)
	}
}
