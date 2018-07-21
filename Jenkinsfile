pipeline {
  agent {
    docker {
        image 'golang:1.9-alpine'
        args '-p 3050:3050 -p 5000:5000'
    }
  }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build') {
      steps {
        echo 'Building...'
        sh 'go version'
      }
    }
    stage('Test') {
      steps {
          sh 'go test ./...'
      }
    }
  } 
}
