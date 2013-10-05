iceray: iceray.go config.go shout.go folderSearch.go songs.go
	go build

dependencies:
	go get

install : iceray
	go install
