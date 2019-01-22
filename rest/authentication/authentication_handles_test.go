package authentication

import (
	"github.com/maxdobeck/gatekeeper/models"
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
		t.Fatal("Problem making our login req. ", err)
	}
	w := httptest.NewRecorder()
	Login(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 401 {
		t.Fail()
	}

	log.Debug(resp.StatusCode)
	log.Debug(resp.Header.Get("Content-Type"))
	log.Debug(string(body))
}

// Test the Login command with a valid set of credentials
func TestLoginGoodCredentials(t *testing.T) {
	models.ConnToDB(os.Getenv("PGURL"))
	_, delErr := models.Db.Query("DELETE FROM members WHERE email like 'validmember@gmail.com'")
	log.Info(delErr)

	valid := models.NewMember{
		Name:      "Valid Test Member",
		Email:     "validmember@gmail.com",
		Email2:    "validmember@gmail.com",
		Password:  "superduper",
		Password2: "superduper",
	}
	models.CreateMember(&valid)

	bodyReader := strings.NewReader(`{"email": "validmember@gmail.com", "password": "superduper"}`)
	req, err := http.NewRequest("POST", "/login", bodyReader)
	if err != nil {
		t.Fatal("Problem logging in ", err)
	}
	w := httptest.NewRecorder()
	Login(w, req)

	resp := w.Result()
	log.Info(resp.StatusCode)

	if resp.StatusCode != 200 {
		t.Fail()
	}
}
