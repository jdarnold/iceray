package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"io"
	
	"github.com/systemfreund/go-libshout"
)

var s shout.Shout
var stream chan <- []byte

func setupShout() {
	mountpoint := IcerayCfg.Server.Mount
	if mountpoint[0] != '/' {
		mountpoint = "/" + mountpoint
	}

	fmt.Printf("Connecting to %s:%d\n",IcerayCfg.Server.Hostname, IcerayCfg.Server.Port)
	
	hostname := flag.String("host", IcerayCfg.Server.Hostname, "shoutcast server name")
	port := flag.Uint("port", IcerayCfg.Server.Port, "shoutcast server source port")
	user := flag.String("user", IcerayCfg.Server.User, "source user name")
	password := flag.String("password", IcerayCfg.Server.Password, "source password")
	mount := flag.String("mountpoint", mountpoint, "mountpoint")

	flag.Parse()

	// Setup libshout parameters
	s = shout.Shout{
		Host:     *hostname,
		Port:     *port,
		User:     *user,
		Password: *password,
		Mount:    *mount,
		Format:   shout.FORMAT_MP3,
		Protocol: shout.PROTOCOL_HTTP,
	}

	// Create a channel where we can send the data
	var err error
	stream, err = s.Open()
	if err != nil {
		log.Fatal("Error opening server " + IcerayCfg.Server.Hostname + " : " + err.Error())
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
	
