pipeline {
  agent {
    docker {
        image 'golang:1.9.7'
        args '-p 3050:3050 -p 5000:5000'
    }
  }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build') {
      steps {
        sh '.jenkins/go-deps.sh'
        sh 'go version'
        sh 'go build'
      }
    }
    stage('Test') {
      steps {
          sh 'go test ./...'
      }
    }
    stage('Deploy to Heroku') {
      when {
        branch 'dev'
      }
    }
    stage('Deploy to Production') {
      when {
        branch 'prod'
      }
    }
  } 
}
