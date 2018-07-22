pipeline {
  agent { dockerfile true }
  environment {
      CI = 'true'
  }
  stages {
    stage('Build') {
      steps {
        sh 'pwd && ls && ls /'
        sh 'ls /root && ls /var/jenkins_home'
        sh 'go version'
        sh './jenkins/scripts/go-build.sh'
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
