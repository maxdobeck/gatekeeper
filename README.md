[![Build Status](https://travis-ci.org/maxdobeck/gatekeeper.svg?branch=dev)](https://travis-ci.org/maxdobeck/gatekeeper)
# Gatekeeper
Service to authenticate users, create sessions, and validate sessions.

## Installing and Running
```
# Starting local server
$ go run main.go

# Run all tests
$ go test ./tests/*

# Compile and create executable
This will place the executable under go/bin
$ go install

# Run the executable
$ gatekeeper

# Useful URLs
/validate - Checks if the supplied session is active
/login - Creates a session and returns an HTTP Only cookie
/logout - Destroys the supplied session
/csrfToken - Upon GET request from client will deliver CSRF token to be used in future requests
```

## Dependencies
	github.com/antonlindstrom/pgstore
	github.com/gorilla/context
	github.com/gorilla/csrf
  github.com/rs/cors
	github.com/urfave/negroni

## Sources
Started with handy helping from:
* [Go Web Examples/Sessions](https://gowebexamples.com/sessions/)
* [Gorilla/Sessions](https://github.com/gorilla/sessions)
