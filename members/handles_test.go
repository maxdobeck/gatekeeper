package members

import (
	"fmt"
	_ "github.com/lib/pq" // github.com/lib/pq
	"github.com/maxdobeck/gatekeeper/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestSignupMemberDuplicateEmail tries to sign up the same email twice
func TestSignupMemberDuplicateEmail(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	_, delErr := models.Db.Query("DELETE FROM members WHERE email like 'testValidCreds@gmail.com'")

	if delErr != nil {
		fmt.Println(delErr)
	}

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
	expectedRes := `{"Status":"Member Not Created","Errors":["Email is already in use."]}`
	if wSignup2.Body.String() != expectedRes {
		t.Error("SignupMember allowed a duplicate email through.", wSignup2.Body)
		t.Fail()
	}
}

// TestChangeMemberEmail attempts to update the user's name value
func TestChangeMemberEmail(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	// Replace this with some sort of DELETE_USER call at some point.  So a cascading del can be performed
	_, delErr := models.Db.Query("DELETE FROM members WHERE email like 'someEmail@gmail.com'")
	fmt.Println(delErr)

	signupBody := strings.NewReader(`{"email": "someEmail@gmail.com", "email2":"someEmail@gmail.com", "password": "supersecret", "password2":"supersecret", "name":"Standard Signup"}`)
	signupReq, signupErr := http.NewRequest("POST", "/members", signupBody)
	if signupErr != nil {
		t.Fail()
	}
	wSignup := httptest.NewRecorder()
	SignupMember(wSignup, signupReq)
	fmt.Println(wSignup)

	var id string
	findErr := models.Db.QueryRow("SELECT id FROM members WHERE email like 'someEmail@gmail.com'").Scan(&id)
	fmt.Println(findErr)

	// Change the email
	body := strings.NewReader(`{"newEmail": "testEmailChange@gmail.com", "newEmail2":"testEmailChange@gmail.com"}`)
	req, err := http.NewRequest("POST", "/members/"+id+"/email", body)
	if err != nil {
		t.Fail()
	}
	wChange := httptest.NewRecorder()
	UpdateMemberEmail(wChange, req)
	fmt.Println(wChange)
}

func TestChangeMemberName(t *testing.T) {
	connStr := os.Getenv("PGURL")
	models.ConnToDB(connStr)
	// Replace this with some sort of DELETE_USER call at some point.  So a cascading del can be performed
	_, delErr := models.Db.Query("DELETE FROM members WHERE email like 'someEmail@gmail.com'")
	fmt.Println(delErr)

	signupBody := strings.NewReader(`{"email": "someEmail@gmail.com", "email2":"someEmail@gmail.com", "password": "supersecret", "password2":"supersecret", "name":"Standard Signup"}`)
	signupReq, signupErr := http.NewRequest("POST", "/members", signupBody)
	if signupErr != nil {
		t.Fail()
	}
	wSignup := httptest.NewRecorder()
	SignupMember(wSignup, signupReq)
	fmt.Println(wSignup)

	// Change the email
	body := strings.NewReader(`{"newEmail": "testEmailChange@gmail.com", "newEmail2":"testEmailChange@gmail.com"}`)
	req, err := http.NewRequest("POST", "/members/", body)
	if err != nil {
		t.Fail()
	}
	wChange := httptest.NewRecorder()
	UpdateMemberEmail(wChange, req)
}
