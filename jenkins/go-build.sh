#!/bin/env sh

sh 'ls /go/bin'
sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
sh 'go build'
sh 'ls /go/bin'