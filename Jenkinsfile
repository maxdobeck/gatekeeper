pipeline {
  agent { dockerfile true }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build for Dev') {
      when {
        branch "dev"
      }
      steps {
        sh 'go version'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && git checkout dev && git pull'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go install'
      }
    }
    stage('Build for Master') {
      when {
        branch "master"
      }
      steps {
        sh 'go version'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && git pull'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go install'
      }
    }
    stage('Build for Prod') {
      when {
        branch "prod"
      }
      steps {
        sh 'go version'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && git checkout prod && git pull'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go install'
      }
    }
    /*stage('Test') {
      steps {
          sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go test ./...'
      }
    }
    stage('Deploy to Dev') {
      when {
        branch 'dev'
      }
      steps {
        echo 'deploying to heroku for dev test'
      }
    }
    stage('Deploy to Prod') {
      when {
        branch 'prod'
      }
      steps { 
        echo 'deploying to heroku prod server'
      }
    }*/
  }
}
