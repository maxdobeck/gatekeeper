pipeline {
  agent { dockerfile { args '-v /go/bin:/go/bin:ro'} }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build') {
      steps {
        sh 'go version'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && git checkout create-jenkinsfile && git pull'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go build && cp gatekeeper /go/bin && ls /go/bin'
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
