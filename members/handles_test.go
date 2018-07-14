package members

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // github.com/lib/pq
	"github.com/maxdobeck/gatekeeper/authentication"
	"github.com/maxdobeck/gatekeeper/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type memberSignup struct {
	Status string   `json:"Status"`
	Errors []string `json:"Errors"`
}

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
	actualRes := memberSignup{}
	json.Unmarshal([]byte(wSignup2.Body.String()), &actualRes)
	var expectedRes [1]string
	expectedRes[0] = "Email is already in use."
	if actualRes.Errors[0] != expectedRes[0] {
		t.Error("SignupMember allowed a duplicate email through.", wSignup2.Body)
	}
}

// TestChangeMemberEmail attempts to update the user's name value
func TestChangeMemberEmail(t *testing.T) {
	// Need to add Mux Router here

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
	if findErr != nil {
		fmt.Println(findErr)
	}

	// Login to start a session
	loginBody := strings.NewReader(`{"email": "someEmail@gmail.com", "password": "supersecret"}`)
	loginReq, loginErr := http.NewRequest("POST", "/login", loginBody)
	if loginErr != nil {
		t.Fail()
	}
	wLogin := httptest.NewRecorder()
	authentication.Login(wLogin, loginReq)

	// Change the email
	body := strings.NewReader(`{"newEmail1": "newEmail@gmail.com", "newEmail2":"newEmail@gmail.com"}`)
	req, err := http.NewRequest("PUT", "/members/"+id+"/email", body)
	if err != nil {
		t.Fail()
	}
	wChange := httptest.NewRecorder()
	req.AddCookie(wLogin.Result().Cookies()[0])

	router := mux.NewRouter()
	router.HandleFunc("/members/{id}/email", UpdateMemberEmail)
	router.ServeHTTP(wChange, req)
	//UpdateMemberEmail(wChange, req)
	fmt.Println("Response:", wChange)
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
