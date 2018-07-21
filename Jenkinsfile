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
        sh 'go version'
        sh 'go get github.com/antonlindstrom/pgstore'
        sh 'go get github.com/gorilla/context'
        sh 'go get github.com/gorilla/csrf'
        sh 'go get github.com/gorilla/mux'
        sh 'go get github.com/lib/pq'
        sh 'go get github.com/rs/cors'
        sh 'go get github.com/urfave/negroni'
        sh 'go build'
      }
    }
    stage('Test') {
      steps {
          sh 'go test ./...'
      }
    }
  } 
}
