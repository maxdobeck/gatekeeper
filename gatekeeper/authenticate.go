package gatekeeper

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkPassword (r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
			panic(err)
	}
	fmt.Println(string(body))
}
