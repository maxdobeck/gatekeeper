#!/bin/sh bash

sh 'ls /go/bin'
sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go build'
sh 'ls/go/bin'