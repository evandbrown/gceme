go clean -i
gox -osarch="linux/amd64" -output="out/{{.OS}}_{{.Arch}}/gceme" 
gsutil cp -a public-read out/linux_amd64/gceme gs://evandbrown17/gceme
