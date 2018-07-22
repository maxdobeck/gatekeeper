#!/bin/env sh

sh 'cd /go/src/github.com/maxdobeck/gatekeeper && git checkout create-jenkinsfile && git pull'
sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go install'
