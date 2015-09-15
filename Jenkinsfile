node('docker') {
  checkout scm

  stage 'test'
  docker.image('golang:1.5.1').inside {
    sh 'go get -d -v'
    sh 'go test'
  }

  stage 'package'
  sh 'gcloud docker -a'
  def img = docker.build('gcr.io/evandbrown17/gceme')
  img.push()
}
