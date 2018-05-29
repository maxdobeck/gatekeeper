package main

import (
	"fmt"
	"github.com/maxdobeck/gatekeeper/authentication"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// An HTTP test to ensure a login request is rejected if the credentials are wrong
func TestInvalidLogin(t *testing.T) {
	bodyReader := strings.NewReader(`{"email": "WrongEmail@email.com", "password": "wrongPassword"}`)

	req, err := http.NewRequest("POST", "/login", bodyReader)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	gatekeeper.Login(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 401 {
		t.Fail()
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}
