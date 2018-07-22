pipeline {
  agent { dockerfile true }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build') {
      steps {
        sh 'go version && echo $GOPATH'
        sh 'go get ./...'
        sh 'go build'
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
