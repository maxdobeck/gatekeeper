pipeline {
  agent { dockerfile true }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build') {
      steps {
        sh 'ls /go/bin'
        sh 'go version'
        // sh 'ls ./jenkins/scripts/go-build.sh'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && git checkout create-jenkinsfile && git pull'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go get ./...'
        sh 'cd /go/src/github.com/maxdobeck/gatekeeper && go install'
        sh 'ls /go/bin'
      }
    }
    stage('Test') {
      steps {
          sh 'go test ./...'
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
    }
  } 
}
