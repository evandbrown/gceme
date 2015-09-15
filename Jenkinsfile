node('docker') {
  checkout scm

  stage 'Go tests'
  docker.image('golang:1.5.1').inside {
    sh('go get -d -v')
    sh('go test')
  }

  stage 'Build Docker image'
  sh('gcloud docker -a')
  def img = docker.build("gcr.io/evandbrown17/gceme:${env.BUILD_TAG}")
  img.push()

  stage 'Deploy to QA cluster'

  stage 'Approve, deploy to prod'
  img.push('latest')
}
