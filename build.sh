go clean -i
gox -osarch="linux/amd64" -output="out/{{.OS}}_{{.Arch}}/gceme" -ldflags "-X github.com/evandbrown/gceme/main.version $(git describe --tags)"
go install -ldflags "-X github.com/evandbrown/gceme/main.version $(git describe --tags)"
gsutil cp -a public-read out/linux_amd64/gceme gs://evandbrown17/gceme
