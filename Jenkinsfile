node('docker') {
  checkout scm
  def zone = sh 'curl -s -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/zone" | grep -o [[:alnum:]-]*$'

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
  docker.image('google/cloud-sdk').inside {
    sh('gcloud components update kubectl --quiet')
    sh('kubectl --namespace=development rollingupdate gceme-frontend --image=${img.id}')
    sh('kubectl --namespace=development rollingupdate gceme-backend --image=${img.id}')
  }

  stage 'Approve, deploy to prod'
  input('Push to prod?')
  img.push('latest')
}
