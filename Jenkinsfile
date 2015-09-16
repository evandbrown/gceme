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
  docker.image('buildpack-deps:jessie-scm').inside {
    sh('export CLOUDSDK_CORE_DISABLE_PROMPTS=1 ; curl https://sdk.cloud.google.com | bash')
    sh('curl -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.0.1/bin/linux/amd64/kubectl')
    sh("/root/google-cloud-sdk/bin/gcloud container clusters get-credentials ${cluster} --zone ${zone}")
    sh("kubectl --namespace=development rollingupdate gceme-frontend --image=${img.id}")
    sh("kubectl --namespace=development rollingupdate gceme-backend --image=${img.id}")
  }

  stage 'Approve, deploy to prod'
  input('Push to prod?')
  img.push('latest')
}
