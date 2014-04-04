package main

import (
	"fmt"
	"log"
	"os"
	"io"
	
	"github.com/systemfreund/go-libshout"
)

var s shout.Shout
var stream chan <- []byte

func setupShout() {
	fmt.Printf("Connecting to %s:%d\n",*Hostname, *Port)
	
	// Setup libshout parameters
	s = shout.Shout{
		Host:     *Hostname,
		Port:     *Port,
		User:     *Username,
		Password: *Password,
		Mount:    *Mount,
		Format:   shout.FORMAT_MP3,
		Protocol: shout.PROTOCOL_HTTP,
	}

	// Create a channel where we can send the data
	var err error
	stream, err = s.Open()
	if err != nil {
		log.Fatal("Error opening server " + *Hostname + " : " + err.Error())
	}
}

func shoutSong( sr SongRecord, buffer []byte ) {
	fd,err := os.Open(sr.fullpath)

	if err != nil {
		log.Println("Problems opening " + sr.fullpath + " : " + err.Error())
		return
	}
	
	// add track info to the stream
	track := sr.title + " by " + sr.artist
	fmt.Println("Playing " + track)
	s.UpdateMetadata( "song", track )
	
	for {
		// Read from file
		n, err := fd.Read(buffer)
		if err != nil && err != io.EOF { panic(err) }
		if n == 0 { break }

		// Send to shoutcast server
		stream <- buffer
	}
}

func closeShout() {
	s.Close()
}
	
