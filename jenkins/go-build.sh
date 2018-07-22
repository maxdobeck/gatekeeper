#!/usr/bin/env bash

sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go build'