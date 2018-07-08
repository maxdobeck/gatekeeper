package members

import (
	"fmt"
	"github.com/maxdobeck/gatekeeper/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestSignupMemberDuplicateEmail tries to sign up the same email twice
func TestSignupMemberDuplicateEmail(t *testing.T) {
	models.ConnToDB(os.Getenv("PGURL"))
	// Signup a user
	signupBody := strings.NewReader(`{"email": "testValidCreds@gmail.com", "email2":"testValidCreds@gmail.com", "password": "supersecret", "password2":"supersecret", "name":"Valid User Signup"}`)
	signupReq, signupErr := http.NewRequest("POST", "/members", signupBody)
	if signupErr != nil {
		t.Fail()
	}
	wSignup := httptest.NewRecorder()
	SignupMember(wSignup, signupReq)

	fmt.Println("Now try and sign up user with same email")

	//Signup same member again
	wSignup2 := httptest.NewRecorder()
	SignupMember(wSignup2, signupReq)
}
