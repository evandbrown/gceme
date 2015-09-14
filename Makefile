gceme: 
	GOOS=linux go build --ldflags="-X main.Version $$USER-$$(TZ=UTC date +%FT%T)Z" -o out/linux/gceme .
	GOOS=darwin go build --ldflags="-X main.Version  $$USER-$$(TZ=UTC date +%FT%T)Z" -o out/darwin/gceme .

# After "make upload", either reboot the machine, or ssh to it and:
#   sudo systemctl restart gobuild.service
# And watch its logs with:
#   sudo journalctl -f -u gobuild.service
upload: gceme
	gsutil cp -a public-read out/linux/gceme gs://evandbrown17/linux/gceme
	gsutil cp -a public-read out/darwin/gceme gs://evandbrown17/darwin/gceme
