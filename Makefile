iceray: iceray.go config.go shout.go
	go build

dependencies:
	go get

install : iceray
	go install
