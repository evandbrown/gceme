node('docker') {
  checkout scm
  def rev = sh 'git rev-parse HEAD'

  stage 'test'
  docker.image('golang:1.5.1').inside('-v /home/jenkins-agent/workspace/$JOB_NAME:/usr/src/JOB_NAME -w /usr/src/JOB_NAME') {
    sh 'go get -d -v'
    sh 'go test'
  }

  stage 'package'
  sh 'gcloud docker -a'
  docker.build('gcr.io/evandbrown17/gceme:latest').push()
}
