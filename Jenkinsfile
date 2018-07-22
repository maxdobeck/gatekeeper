pipeline {
  agent {
    docker {
        image 'golang:1.9.7'
        args 'ADD /var/jenkins_home/.netrc /root/.netrc'
        args '-p 3050:3050 -p 5000:5000'
    }
  }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build') {
      steps {
        sh 'go get ./...'
        sh 'cat .netrc'
        sh 'go version'
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
