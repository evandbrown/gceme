node('docker') {
  checkout scm
  def go = docker.image 'golang:1.5.1'
  stage 'test'
  go.inside "-v /home/jenkins-agent/workspace/$JOB_NAME:/usr/src/$JOB_NAME -w /usr/src/$JOB_NAME" {
    sh 'go get -d -v'
    sh 'go test'
  }

  stage 'package'
  sh 'gcloud docker -a'
  def img = go.build "gcr.io/evandbrown17/gceme:${env.BUILD_TAG}"
  img.push()
  img.push 'latest'
}
