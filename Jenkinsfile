node('docker') {
  checkout scm

  // Kubernetes cluster info
  def cluster = 'jenkins'
  def zone = sh 'curl -s -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/zone" | grep -o [[:alnum:]-]*$'

  // Run tests
  stage 'Go tests'
  docker.image('golang:1.5.1').inside {
    sh('go get -d -v')
    sh('go test')
  }

  // Build image with Go binary
  stage 'Build Docker image'
  sh('gcloud docker -a')
  def img = docker.build("gcr.io/evandbrown17/gceme:${env.BUILD_TAG}")
  img.push()

  // Deploy image to cluster in dev namespace
  stage 'Deploy to QA cluster'
  docker.image('google/cloud-sdk').inside {
    sh('gcloud components update kubectl --quiet')
    sh("gcloud container clusters get-credentials ${cluster} --zone ${zone}")
    sh('kubectl --namespace=development rollingupdate gceme-frontend --image=${img.id}')
    sh('kubectl --namespace=development rollingupdate gceme-backend --image=${img.id}')
  }

  stage 'Approve, deploy to prod'
  input('Push to prod?')
  img.push('latest')
}
