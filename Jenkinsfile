node('docker') {
  docker.image('golang:1.5.1').inside('-v /home/jenkins-agent/workspace/$JOB_NAME:/usr/src/JOB_NAME -w /usr/src/JOB_NAME') {
    git 'https://github.com/evandbrown/gceme.git'
    sh 'go get -d -v'
    sh 'go test'
  }
}
