node('docker') {
  checkout scm

  // Kubernetes cluster info
  def cluster = 'jenkins'
  def zone = 'us-central1-f'

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
    sh('apt-get update -y ; apt-get install -y curl')
    sh('curl -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kubectl')
    sh("gcloud container clusters get-credentials ${cluster} --zone ${zone}")
    sh("kubectl --namespace=development rollingupdate gceme-frontend --image=${img.id}")
    sh("kubectl --namespace=development rollingupdate gceme-backend --image=${img.id}")
  }

  stage 'Approve, deploy to prod'
  input('Push to prod?')
  img.push('latest')
}
