FROM golang:1.9.7

ENV GOBIN=$GOPATH/bin

# Must use an env script w/ an env variable to hide the api token if this is ever shared.
RUN git config --global url."https://b99b85d769f70db60e133e6dbcefb83c877b8f01@github.com/".insteadOf "https://github.com/"

RUN mkdir -p src/github.com/maxdobeck

WORKDIR /go/src/github.com/maxdobeck/

RUN git clone https://github.com/maxdobeck/gatekeeper.git

WORKDIR /go/src/github.com/maxdobeck/gatekeeper

# This was to ensure that the git checkout process worked and and the go get cmd worked.
#RUN git branch -a

# RUN git checkout create-jenkinsfile

# RUN go get ./...

# RUN go build

# This is for actually running the binary.
# TODO: Add env variables.
# TODO: List env variables for record keeping
# RUN /go/bin/gatekeeper