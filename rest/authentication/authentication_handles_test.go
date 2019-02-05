package authentication

import (
	"github.com/maxdobeck/gatekeeper/models"
	"github.com/maxdobeck/gatekeeper/rest/members"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// An HTTP test to ensure a login request is rejected if the credentials are wrong
func TestLoginBadCredentials(t *testing.T) {
	bodyReader := strings.NewReader(`{"email": "WrongEmail@email.com", "password": "wrongPassword"}`)

	req, err := http.NewRequest("POST", "/login", bodyReader)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	Login(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 401 {
		t.Fail()
	}

	log.Info(resp.StatusCode)
	log.Info(resp.Header.Get("Content-Type"))
	log.Info(string(body))
}

// Test the Login command with a valid set of credentials
func TestLoginGoodCredentials(t *testing.T) {
	models.ConnToDB(os.Getenv("PGURL"))
	// Signup a user
	signupBody := strings.NewReader(`{"email": "testValidCreds@gmail.com", "email2":"testValidCreds@gmail.com", "password": "supersecret", "password2":"supersecret", "name":"Valid User Signup"}`)
	signupReq, signupErr := http.NewRequest("POST", "/members", signupBody)
	if signupErr != nil {
		t.Log("Ignoring that the user already exists.  Doesn't matter for test.")
	}
	wSignup := httptest.NewRecorder()
	members.SignupMember(wSignup, signupReq)

	bodyReader := strings.NewReader(`{"email": "testValidCreds@gmail.com", "password": "supersecret"}`)
	req, err := http.NewRequest("POST", "/login", bodyReader)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	Login(w, req)

	resp := w.Result()
	log.Info(resp.StatusCode)

	if resp.StatusCode != 200 {
		t.Fail()
	}
}
