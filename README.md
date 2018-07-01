[![Build Status](https://travis-ci.org/maxdobeck/gatekeeper.svg?branch=dev)](https://travis-ci.org/maxdobeck/gatekeeper)
# Gatekeeper
Service to authenticate users, manage user sessions, and process business logic.

## Installing and Running
```
# Starting local server
$ go run main.go

# Run all tests
$ go test ./...

# Compile and create executable
This will place the executable under go/bin
$ go install

# Run the executable
$ gatekeeper
```

## Dependencies
```
	github.com/antonlindstrom/pgstore
	github.com/gorilla/context
	github.com/gorilla/csrf
  github.com/rs/cors
	github.com/urfave/negroni
```

## Sources
Started with handy helping from:
* [Go Web Examples/Sessions](https://gowebexamples.com/sessions/)
* [Gorilla/Sessions](https://github.com/gorilla/sessions)
