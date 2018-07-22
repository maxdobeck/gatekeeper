#!/bin/env sh

sh 'ls go/bin'
sh '/go/src/github.com/maxdobeck/gatekeeper/git checkout create-jenkinsfile && /go/src/github.com/maxdobeck/gatekeeper/git pull'
sh '/go/src/github.com/maxdobeck/gatekeeper/go get ./...'
sh '/go/src/github.com/maxdobeck/gatekeeper/go build'
sh 'ls go/bin'
