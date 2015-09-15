node('docker') {
  checkout scm
  def go = docker.image 'golang:1.5.1'
  stage 'test'
  go.inside {
    sh 'go get -d -v'
    sh 'go test'
  }

  stage 'package'
  sh 'gcloud docker -a'
  def img = go.build "gcr.io/evandbrown17/gceme:latest"
  img.push()
  img.push 'latest'
}
