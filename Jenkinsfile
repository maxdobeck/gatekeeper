pipeline {
  agent { dockerfile { args '--net="host" -e PGURL=postgres://sfdev:sfdev@localhost:5432/scheduler_fairy_dev -e GO_ENV=dev -p 3050:3050 -e PORT=3050' } }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build for Dev') {
      when {
        branch "dev"
      }
      steps {
        sh 'go version && git branch'
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
        sh 'go version && git branch'
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
        sh 'go version && git branch'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && git checkout prod && git pull'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go install'
      }
    }
    stage('Test') {
      steps {
        sh 'git branch'
          sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go test ./...'
      }
    }
    stage('Deploy to Dev') {
      when {
        branch 'dev'
      }
      steps {
        echo 'Ready for deploy to Heroku'
      }
    }
    /*stage('Deploy to Prod') {
      when {
        branch 'prod'
      }
      steps { 
        echo 'deploying to heroku prod server'
      }
    }*/
  }
}
