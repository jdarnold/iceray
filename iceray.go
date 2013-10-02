// iceray project main.go
package main

import (
	"fmt"
	"os"
	"flag"
	"io"
	"path"
	"strings"
	"log"
	"github.com/jdarnold/go-id3/src/id3"
	"github.com/jdarnold/go-libshout"
)

type SongRecord struct {
	fullpath string
	filetype string
	title string
	artist string
}

// Setup some command line flags
var (
    hostname = flag.String("host", "amazingdev.com", "Vehement Flame Radio")
    port = flag.Uint("port", 8000, "shoutcast server source port")
    user = flag.String("user", "source", "source user name")
    password = flag.String("password", "flamecast", "source password")
    mount = flag.String("mountpoint", "/flame", "mountpoint")
) 

func sdir(folder string, addfilechannel chan SongRecord) {
	searchdir, eopen := os.Open(folder)
	if eopen != nil {
		log.Println("Error opening " + folder + " : " + eopen.Error())
		return
	}
	
	homefiles, eread := searchdir.Readdir(-1)
	if eread != nil {
		log.Println("Error reading " + folder + " : " + eopen.Error())
		return
	}
	for i := range homefiles {
		fname := homefiles[i].Name()
		if fname[0] == '.' {
			continue
		}

		if homefiles[i].IsDir() {
			ndir := folder+"/"+fname
			fmt.Println("Searching " + ndir)
			go sdir(ndir,addfilechannel)
			continue
		}

		if !strings.Contains(fname,".mp3") {
			continue
		}
			
		if homefiles[i].Size() < 100 {
			continue
		}

		var fext = path.Ext(fname);

		var sr SongRecord
		sr.fullpath = folder+"/"+fname
		sr.filetype = strings.TrimPrefix(fext,".")
		sr.title = "Unknown"
		sr.artist = "Unknown"
		
		addfilechannel <- sr
	}
}

func main() {
	folder := "."

	if len(os.Args) > 1 {
		folder = os.Args[1]
	}
	
	addfilechannel := make(chan SongRecord, 10)
	defer close(addfilechannel)

	go sdir(folder,addfilechannel)

	flag.Parse()

	// Setup libshout parameters
	s := shout.Shout{
		Host:     *hostname,
		Port:     *port,
		User:     *user,
		Password: *password,
		Mount:    *mount,
		Format:   shout.FORMAT_MP3,
		Protocol: shout.PROTOCOL_HTTP,
	}

	defer s.Close()

	// Create a channel where we can send the data
	//
	stream, err := s.Open()
	if err != nil {
		panic(err)
	}
	
	buffer := make([]byte, shout.BUFFER_SIZE)
	
	for {
		mfile := <- addfilechannel
		fd,err := os.Open(mfile.fullpath)
		defer fd.Close()
		
		if err != nil {
			log.Println("Problem opening: " + mfile.fullpath)
			continue
		}
		
		mp3tags := id3.Read(fd)

		if ( mp3tags != nil ) {
			mfile.artist = mp3tags.Artist
		}

		if ( mp3tags != nil ) {
			mfile.title = mp3tags.Name
		}

		fmt.Println("Playing " + mfile.title + " by " + mfile.artist )

		fd.Seek(0,0)

		
		s.UpdateMetadata( "song", mfile.title + " by " + mfile.artist )
		
		for {
			// Read from file
			n, err := fd.Read(buffer)
			if err != nil && err != io.EOF { panic(err) }
			if n == 0 { break }

			// Send to shoutcast server
			stream <- buffer
		}

	}
}

















