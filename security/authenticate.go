package gatekeeper

import (
	"fmt"
	"encoding/json"
	"net/http"
)

// Credentials are the user provided email and password
type Credentials struct {
	Email, Password string
}

func checkPassword(r *http.Request) {
	var c Credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	fmt.Println(c.Email, c.Password)
	if err != nil {
		fmt.Println("Error decoding credentials >>", err)
	}
}
