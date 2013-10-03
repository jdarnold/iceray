iceray: iceray.go
	go build

dependencies:
	go get

install : iceray
	go install
