pipeline {
  agent { dockerfile true }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build') {
      steps {
        sh 'go version'
        sh 'chmod 0755 jenkins/scripts/go-build.sh'
        sh 'ls -l ./jenkins/scripts/go-build.sh'
        sh 'pwd && ls . && jenkins/scripts/go-build.sh'
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
